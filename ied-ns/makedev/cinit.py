import os
import subprocess

import jinja2 as j2


class CINITResult:
    def __init__(self, iso_p: str, user_data: str, cloud_data: str):
        self.iso_p = iso_p
        self.user_data = user_data
        self.cloud_data = cloud_data

    def __str__(self):
        return f"ISO: {self.iso_p}, User Data: {self.user_data}, Cloud Data: {self.cloud_data}"


class FileWrite:
    def __init__(self, path: str, owner: str, permissions: str, content: str):
        self.path = path
        self.owner = owner
        self.permissions = permissions
        self.content = content


class UserData:
    def __init__(
        self,
        hostname: str,
        password: str,
        commands: list[str],
        writes: list[FileWrite],
    ):
        self.hostname = hostname
        self.password = password
        self.commands = commands
        self.writes = writes


def prepare(
    devname: str,
    cmds: list[str],
    flws: list[FileWrite],
    fp: str = "",
    write: bool = False,
) -> CINITResult:
    jenv = j2.Environment(
        loader=j2.PackageLoader("makedev"), trim_blocks=True, lstrip_blocks=True
    )
    jtempl = jenv.get_template("user-data.jinja")
    out = jtempl.render(
        ud=UserData(
            hostname=devname,
            password="root",  # is this fine???
            commands=cmds,
            writes=flws,
        )
    )
    if write:
        os.makedirs(fp, exist_ok=True)
        os.makedirs(f"{fp}/seed", exist_ok=True)
        with open(f"{fp}/seed/user-data", "w") as f:
            f.write(out)
        with open(f"{fp}/seed/meta-data", "w") as f:
            f.write("")
        subprocess.run(
            [
                "genisoimage",
                "-o",
                os.path.abspath(f"{fp}/cloudinit.iso"),
                "-volid",
                "cidata",
                "-joliet",
                "-rock",
                f"{fp}/seed/user-data",
                f"{fp}/seed/meta-data",
            ],
            check=True,
        )
    return CINITResult(
        iso_p=f"{fp}/cloudinit.iso",
        user_data=f"{fp}/user-data",
        cloud_data=f"{fp}/cloud-data",
    )

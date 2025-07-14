import os
import subprocess

import jinja2 as j2


class FileWrite:
    def __init__(self, path, owner, permissions, content):
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


def write(ud: UserData, fp: str | None = None) -> tuple[str, str, str]:
    jenv = j2.Environment(
        loader=j2.PackageLoader("cloud-init"), trim_blocks=True, lstrip_blocks=True
    )
    jtempl = jenv.get_template("user-data.jinja")
    out = jtempl.render(ud=ud)
    if fp is None:
        print(f"Cloud-init configuration for {ud.hostname}")
        print(out)
        return ("", "", "")
    else:
        os.makedirs(fp, exist_ok=True)
        with open(f"{fp}/user-data", "w") as f:
            f.write(out)
        with open(f"{fp}/cloud-data", "w") as f:
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
                f"{fp}/user-data",
                f"{fp}/meta-data",
            ],
            check=True,
        )
        return (f"{fp}/cloudinit.iso", f"{fp}/user-data", f"{fp}/cloud-data")


# testing only!!!!
if __name__ == "__main__":
    write(
        UserData(
            hostname="hostname",
            password="password",
            commands=["c", "o"],
            writes=[
                FileWrite("/etc/hostname", "root", "0644", "hostname\n"),
                FileWrite("/etc/hosts", "root", "0644", "hosts\n"),
            ],
        )
    )

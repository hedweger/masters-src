import os
import subprocess

import jinja2 as j2


class NetworkConnection:
    def __init__(
        self, name: str, iface: str, src_ip: str, mac: str, gateway: str | None
    ):
        self.name: str = name
        self.iface: str = iface
        self.src_ip: str = src_ip
        self.mac: str = mac
        self.gateway: str | None = gateway


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
        dev_type: str,
        hostname: str,
        password: str,
        commands: list[str],
        writes: list[FileWrite],
    ):
        self.dev_type = dev_type
        self.hostname = hostname
        self.password = password
        self.commands = commands
        self.writes = writes


def prepare(
    dev_type: str,
    devname: str,
    cmds: list[str],
    flws: list[FileWrite],
    connections: list[NetworkConnection] = [],
    fp: str = "",
    write: bool = False,
) -> CINITResult:
    jenv = j2.Environment(
        loader=j2.PackageLoader("makedev"), trim_blocks=True, lstrip_blocks=True
    )
    jtempl = jenv.get_template("user-data.jinja")
    out = jtempl.render(
        ud=UserData(
            dev_type=dev_type,
            hostname=devname,
            password="root",  # is this fine???
            commands=cmds,
            writes=flws,
        )
    )
    jtempl = jenv.get_template("network-config.jinja")
    # if dev_type == "switch":
    #     connections.append(
    #         NetworkConnection(
    #             name=f"br0",
    #             iface=f"br0",
    #             src_ip='',
    #             mac='',
    #             gateway='',
    #         )
    #     )
    #
    netw_out = jtempl.render(connections=connections, dev_type=dev_type)
    if write:
        os.makedirs(fp, exist_ok=True)
        os.makedirs(f"{fp}/seed", exist_ok=True)
        with open(f"{fp}/seed/user-data", "w") as f:
            f.write(out)
        with open(f"{fp}/seed/network-config", "w") as f:
            f.write(netw_out)
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
                f"{fp}/seed/network-config",
            ],
            check=True,
        )
    return CINITResult(
        iso_p=f"{fp}/cloudinit.iso",
        user_data=f"{fp}/user-data",
        cloud_data=f"{fp}/cloud-data",
    )

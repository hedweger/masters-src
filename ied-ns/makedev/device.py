from enum import Enum

import jinja2 as j2
from makedev.cinit import FileWrite, NetworkConnection


class DeviceType(Enum):
    SW = "switch"
    RTU = "rtu"


class Device:
    def __init__(
        self,
        dev_type: DeviceType,
        name: str,
        address: str,
    ):
        self.dev_type: DeviceType = dev_type
        self.networks: list[NetworkConnection] = []

        self.name: str = name
        self.address: str = address if "/" in address else f"{address}/24"
        self.iface_count = 0

        self.image_path: str | None = None
        self.seed_path: str | None = None
        self.user_data_path: str | None = None
        self.cloud_data_path: str | None = None

    def add_network_connection(self, network_name: str, mac: str, gateway: str | None):
        self.iface_count += 1
        self.networks.append(
            NetworkConnection(
                name=network_name,
                iface=f"ens{self.iface_count + 1}",
                src_ip=self.address,
                mac=mac,
                gateway=gateway,
            )
        )

    def startup_commands(self) -> list[str]:
        if self.dev_type is DeviceType.SW:
            return [
                "ip link add name br0 type bridge",
                "ip link set dev ens2 master br0",
                "ip link set dev ens3 master br0",
                "ip link set dev br0 up",
                "sudo tcpdump -i br0 not arp and not llc",
            ]
        else:
            cmds = []
            cmds.append(
                "sudo wget https://github.com/hedweger/masters-src/releases/download/client/ied-client"
            )
            cmds.append("sudo chmod +x /ied-client")
            if self.name == "pc1":
                cmds.append(
                    "sudo wget https://github.com/hedweger/masters-src/releases/download/server/ied-server.tar"
                )
                cmds.append("sudo tar -xf /ied-server.tar")
                cmds.append("sudo chmod +x /ied-server")
                cmds.append("sudo /ied-server")
            return cmds

    def startup_filewrites(self) -> list[FileWrite]:
        if self.dev_type is DeviceType.SW:
            return []
        else:
            return []

    def libvirt_xml(self) -> str:
        jenv = j2.Environment(
            loader=j2.PackageLoader("makedev"), trim_blocks=True, lstrip_blocks=True
        )
        jtempl = jenv.get_template("virt_device.xml.jinja")
        return jtempl.render(
            dtype=self.dev_type.value,
            name=self.name,
            nram="512",  # for now
            vcpu="1",  # for now
            disk=self.image_path,
            seed=self.seed_path,
            nets=self.networks,
        )

from enum import Enum

import jinja2 as j2
from makedev.cinit import FileWrite


class NetworkConnection:
    def __init__(self, name: str, src_ip: str, mac: str, gateway: str):
        self.name: str = name
        self.src_ip: str = src_ip
        self.mac: str = mac
        self.gateway: str = gateway


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

        self.image_path: str | None = None
        self.seed_path: str | None = None
        self.user_data_path: str | None = None
        self.cloud_data_path: str | None = None

    def add_network_connection(self, network_name: str, mac: str, gateway: str):
        self.networks.append(
            NetworkConnection(
                name=network_name,
                src_ip=self.address,
                mac=mac,
                gateway=gateway,
            )
        )

    def startup_commands(self) -> list[str]:
        if self.dev_type is DeviceType.RTU:
            return [
                "sleep 5",
                "sudo systemctl restart systemd-networkd",
                "sudo ip link set dev ens3 up",
                "stty erase ^H",
            ]
        elif self.dev_type is DeviceType.SW:
            return [
                # "ip link add name br0 type bridge",
                # "ip link set dev ens2 master br0",
                # "ip link set dev ens3 master br0",
                # "ip link set dev br0 up",
                # "ip link set dev ens2 up",
                # "ip link set dev ens3 up",
                "sudo systemctl restart systemd-networkd",
                "stty erase ^H",
                # "ip link set dev ens4 up",
            ]
        return []

    def startup_filewrites(self) -> list[FileWrite]:
        if self.dev_type is DeviceType.SW:
            return (
                [
                    FileWrite(
                        path=f"/etc/systemd/network/0{i+3}-ens{i+2}.network",
                        content=f"[Match]\nName=ens{i+2}\n\n[Network]\nBridge=br0\n",
                        owner="root:root",
                        permissions="0644",
                    )
                    for i, netw in enumerate(self.networks)
                ]
                + [
                    FileWrite(
                        path=f"/etc/systemd/network/01-ens{len(self.networks)+2}.network",
                        content=f"[Match]\nName=ens{len(self.networks)+2}\n\n[Network]\nAddress=192.168.122.4/24\nGateway=192.168.122.1\n",
                        owner="root:root",
                        permissions="0644",
                    )
                ]
                + [
                    FileWrite(
                        path=f"/etc/systemd/network/01-br0.netdev",
                        content=f"[NetDev]\nName=br0\nKind=bridge\n\n[Bridge]\nSTP=on\n",
                        owner="root:root",
                        permissions="0644",
                    )
                ]
                + [
                    FileWrite(
                        path=f"/etc/systemd/network/02-br0.network",
                        content=f"[Match]\nName=br0\n\n[Network]\n",
                        owner="root:root",
                        permissions="0644",
                    )
                ]
            )
        elif self.dev_type is DeviceType.RTU:
            return [
                FileWrite(
                    path=f"/etc/systemd/network/01-ens{i+2}.network",
                    content=f"[Match]\nName=ens{i+2}\nType=ether\n\n[Network]\nAddress={netw.src_ip}\n",
                    owner="root:root",
                    permissions="0644",
                )
                for i, netw in enumerate(self.networks)
            ] + [
                FileWrite(
                    path=f"/etc/systemd/network/01-ens{len(self.networks)+2}.network",
                    content=f"[Match]\nName=ens{len(self.networks)+2}\n\n[Network]\nAddress=192.168.122.{int(self.name[-1])+1}/24\nGateway=192.168.122.1\n",
                    owner="root:root",
                    permissions="0644",
                )
            ]
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

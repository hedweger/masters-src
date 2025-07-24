import ipaddress as ip
import os
from typing import Iterable

import jinja2 as j2
import makedev.cinit as cinit
import makedev.drive as drive
import makedev.virtmac as mac
import yaml
from makedev.device import Device, DeviceType


class DeviceManager:
    def __init__(self):
        self.devices: dict[str, Device] = {}
        self.context: str
        self.network_address: ip.IPv4Network
        self.mac_iter = mac.gen()
        self.ip4_iter: Iterable

    def parse(self, cfg_fp: str, write: bool = False):
        with open(cfg_fp, "r") as file:
            config_data = yaml.safe_load(file)
        # create file path context for later use
        self.context = "/".join(cfg_fp.split("/")[:-1]) + "/tmp"
        self.network_address = ip.IPv4Network(config_data["network"]["address"])
        self.ip4_iter = self.network_address.hosts()
        rtus = config_data["rtus"]
        for i, rtu in enumerate(rtus):
            name = rtu["name"] if rtu["name"] else f"rtu{i}"
            address = rtu["address"] if rtu["address"] else f"{next(self.ip4_iter)}/24"
            self.devices[name] = Device(
                DeviceType.RTU,
                name,
                address,
            )

        switches = config_data["switches"]
        for i, sw in enumerate(switches):
            name = sw["name"] if sw["name"] else f"switch{i}"
            address = sw["address"] if sw["address"] else f"{next(self.ip4_iter)}/24"
            self.devices[name] = Device(
                DeviceType.SW,
                name,
                address,
            )
            for conn in sw["connected"]:
                self.create_network(src_name=name, dst_name=conn["to"], gtw_addr="")
        self.create_devices()

    def create_devices(self):
        for device in self.devices.values():
            dev_addr = device.address
            device.address = f'192.168.122.{int(device.name[-1])+1}/24'
            device.add_network_connection(
                network_name=f"default",
                mac=next(self.mac_iter),
                gateway="192.168.122.1",
            )
            device.address = dev_addr
            device.image_path = drive.qcow2(
                f"{self.context}/{device.name}", device.name, True
            )
            seeds = cinit.prepare(
                device.dev_type.value,
                device.name,
                device.startup_commands(),
                device.startup_filewrites(),
                device.networks,
                f"{self.context}/{device.name}",
                True,
            )
            device.seed_path = seeds.iso_p
            device.user_data_path = seeds.user_data
            device.cloud_data_path = seeds.cloud_data
            with open(f"{self.context}/{device.name}/config.xml", "w") as f:
                f.write(device.libvirt_xml())

    def create_network(
        self,
        src_name: str,
        dst_name: str,
        gtw_addr: str,
    ):
        dst_device = self.devices[dst_name]
        src_device = self.devices[src_name]
        gateway = gtw_addr if "/" in gtw_addr else f"{gtw_addr}/24"
        network_name = f"{src_name}-{dst_name}"
        dst_device.add_network_connection(
            network_name=network_name,
            gateway=None,
            mac=next(self.mac_iter),
        )
        src_device.add_network_connection(
            network_name=network_name,
            gateway=None,
            mac=next(self.mac_iter),
        )
        jenv = j2.Environment(
            loader=j2.PackageLoader("makedev"), trim_blocks=True, lstrip_blocks=True
        )
        jtempl = jenv.get_template("virt_network.xml.jinja")
        os.makedirs(self.context, exist_ok=True)
        with open(f"{self.context}/{network_name}.xml", "w") as f:
            f.write(
                jtempl.render(
                    name=network_name,
                    brdg=f"{network_name}-br",
                )
            )

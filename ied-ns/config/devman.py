import ipaddress as ip
import logging as log

import config.cinit as cinit
import config.drive as drive
import yaml
from config.device import Device, DeviceType


class DeviceManager:
    def __init__(self, iface: ip.IPv4Network | None = None):
        self.devs: list[Device] = []
        self.conns: list[tuple[Device, Device]] = []
        if iface is None:
            iface = ip.IPv4Network("192.168.0.0/16")
        self.iface = iface
        self.context: str

    def parse(self, cfg_fp: str):
        logger = log.getLogger(__name__)
        logger.info(f"began parsing configuration file: {cfg_fp}")
        self.context = "/".join(cfg_fp.split("/")[:-1])
        with open(cfg_fp, "r") as file:
            config_data = yaml.safe_load(file)
        devices = config_data.get("devices", [])
        if devices is None:
            raise ValueError("No devices found in the configuration file.")
        for device in devices:
            dtype = device.get("type")
            name = device.get("name")
            addr = device.get("addr")
            conns = device.get("connected_to", [])
            self.add_device(
                dtype,
                name,
                ip.IPv4Address(addr) if addr else None,
                conns,
            )

    def add_device(
        self,
        dev_type: DeviceType,
        name: str,
        addr: ip.IPv4Address | None = None,
        conn: list | None = None,
    ):
        if addr is None:
            addr = self.iface[len(self.devs) + 1]
        if addr not in self.iface:
            raise ValueError(f"Address {addr} is not in the interface {self.iface}")
        self.devs.append(Device(dev_type, name, addr, conn or []))
        log.getLogger(__name__).info(
            f"Added device: {dev_type} {name} with address {addr}"
        )

    def prepare_devs(self):
        for dev in self.devs:
            image = drive.qcow2(self.context, dev.name)
            seeds = cinit.write(
                cinit.UserData(
                    hostname=dev.name,
                    password="root",
                    commands=[],
                    writes=[],
                ),
                f"{self.context}/{dev.name}",
            )
            print(image)
            print(seeds)

    def list(self):
        for dev in self.devs:
            print(dev.__repr__())

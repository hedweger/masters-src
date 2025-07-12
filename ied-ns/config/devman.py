from config.device import Device, DeviceType
import ipaddress as ip
import yaml


class DeviceManager:
    def __init__(self, iface: ip.IPv4Network | None = None):
        self.devs: list[Device] = []
        self.conns: list[tuple[Device, Device]] = []
        if iface is None:
            iface = ip.IPv4Network("192.168.0.0/16")
        self.iface = iface

    def parse(self, cfg_fp: str):
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

    def list(self):
        for dev in self.devs:
            print(dev.__repr__())

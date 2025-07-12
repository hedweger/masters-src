from enum import Enum
import ipaddress as ip


class DeviceType(Enum):
    SW = "switch"
    RTU = "rtu"


class Device:
    def __init__(self, dtype: DeviceType, name: str, addr: ip.IPv4Address, conn: list):
        self.dtype: DeviceType = dtype
        self.name: str = name
        self.addr: ip.IPv4Address = addr
        self.conn: list = conn

    def __repr__(self):
        return f"{self.dtype}(\n\t - name={self.name}, \n\t - addr={self.addr}, \n\t - connected_to={self.conn}\n)"

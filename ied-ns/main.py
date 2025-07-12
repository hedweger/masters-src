import yaml
from config.devman import DeviceManager
import ipaddress as ip


def main(config: str):
    dman = DeviceManager()
    dman.parse(config)
    dman.list()

if __name__ == "__main__":
    main("config.yaml")

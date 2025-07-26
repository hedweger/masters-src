import logging
import os

from makedev.devman import DeviceManager

logger = logging.getLogger(__name__)


def main(config: str):
    dman = DeviceManager()
    dman.parse(config, write=True)


if __name__ == "__main__":
    main(os.path.abspath("../config.yaml"))

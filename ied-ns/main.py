import logging, os
from config.devman import DeviceManager

logger = logging.getLogger(__name__)

def main(config: str):
    dman = DeviceManager()
    dman.parse(config)
    dman.prepare_devs()


if __name__ == "__main__":
    main(os.path.abspath("./config.yaml"))

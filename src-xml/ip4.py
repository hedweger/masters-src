class IPv4:
    def __init__(self, addr: str, prefix: int):
        parts = addr.split(".")
        self.first = parts[1]
        self.second = parts[2]
        self.third = parts[3]
        self.fourth = parts[4]
        self.prefix = prefix

    def addr(self):
        return f"{self.first}.{self.second}.{self.third}.{self.fourth}"

    def netmask(self):
        mask32 = (0xFFFFFFFF << (32 - self.prefix)) & 0xFFFFFFFF
        return ".".join(str((mask32 >> (24 - 8 * i)) & 0xFF) for i in range(4))

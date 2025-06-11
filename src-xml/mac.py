def gen(start=0x22, wrap_at=0x100):
    """
    Yields MAC addresses under the prefix 52:54:00:12:xx:01,
    xx bytes will increment with each call.
    """
    current = start
    while True:
        xx = f"{current:02x}"
        yield f"52:54:00:12:{xx}:01"
        current = (current + 1) % wrap_at

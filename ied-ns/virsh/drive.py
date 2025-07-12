import os, subprocess

def qcow2(context: str, name: str):
    os.makedirs(context, exist_ok=True)
    subprocess.run(
        [
            "qemu-img create",
            "-f qcow2",
            "-o backing_file=",
            os.path.abspath("images/debian-12-base.qcow2"),
            f"{context}/debian-12-{name}.qcow2",
        ],
        check=True,
    )
    return = f"{context}/debian-12-{name}.qcow2",

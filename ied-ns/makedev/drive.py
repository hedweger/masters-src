import os, subprocess

def qcow2(context: str, name: str, write: bool = False) -> str:
    if write:
        os.makedirs(context, exist_ok=True)
        subprocess.run(
            [
                "cp",
                os.path.abspath("debian-12-genericcloud-amd64.qcow2"),
                f"{context}/debian-12-{name}.qcow2",
            ],
            check=True,
        )
    return f"{context}/debian-12-{name}.qcow2"

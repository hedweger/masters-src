import os, subprocess


def prep(
    logger,
    fp: str,
    name: str,
):
    logger.info(f"preparing qcow2 files with fp: {fp}")
    os.makedirs(f"{fp}", exist_ok=True)
    subprocess.run(
        [
            "cp",
            os.path.abspath("images/debian-12-base.qcow2"),
            f"{fp}/debian-12-{name}.qcow2",
        ],
        check=True,
    )
    return f"{fp}/debian-12-{name}.qcow2"

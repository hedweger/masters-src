import os, subprocess


def iso(
    logger,
    fp: str,
    name: str,
    pkgs: list | None,
    wfs: list | None,
    cmds: list | None,
    hostname="localhost",
    pword="root",
):
    logger.info("begin ISO define cloud-install")
    template = f"""#cloud-config\nhostname: {hostname}\npassword: {pword}\nchpasswd: {{ expire: False }}\nssh_pwauth: True
    """

    if pkgs:
        template += "\npackages:\n"
        for pkg in pkgs:
            template += "  - " + pkg

    if cmds:
        template += "\nruncmd:\n"
        for cmd in cmds:
            template += "  - " + cmd + "\n"

    if wfs:
        template += "\nwrite_files:\n"
        for wf in wfs:
            path, content = wf[0], wf[1]
            template += f"  - path: {path}\n"
            template += "    owner: root:root\n"
            template += "    permissions: '0644'\n"
            template += "    content: |\n"
            for line in content.splitlines():
                template += f"      {line}\n"

    os.makedirs(f"{fp}/seed", exist_ok=True)
    with open(f"{fp}/seed/user-data", "w") as f:
        f.write(template)
        logger.info(f"succesfully written user-data at {fp}")

    with open(f"{fp}/seed/meta-data", "w") as f:
        f.write("")
        logger.info(f"succesfully written meta-data at {fp}")

    subprocess.run(
        [
            "genisoimage",
            "-o",
            os.path.abspath(f"{fp}/seed-{name}.iso"),
            "-volid",
            "cidata",
            "-joliet",
            "-rock",
            f"{fp}/seed/user-data",
            f"{fp}/seed/meta-data",
        ],
        check=True,
    )
    logger.info(f"succesfully written seed-{name}.iso at {fp}")
    return f"{fp}/seed-{name}.iso"

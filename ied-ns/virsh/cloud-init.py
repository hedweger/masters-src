import jinja2 as j2


class FileWrite:
    def __init__(self, path, owner, permissions, content):
        self.path = path
        self.owner = owner
        self.permissions = permissions
        self.content = content


class UserData:
    def __init__(
        self,
        hostname: str,
        password: str,
        commands: list[str],
        writes: list[FileWrite],
    ):
        self.hostname = hostname
        self.password = password
        self.commands = commands
        self.writes = writes


def write(ud: UserData, fp: str | None = None):
    jenv = j2.Environment(
        loader=j2.PackageLoader("cloud-init"), trim_blocks=True, lstrip_blocks=True
    )
    jtempl = jenv.get_template("user-data.jinja")
    out = jtempl.render(ud=ud)
    if fp is None:
        print(out)
    else:
        with open(f"{fp}/user-data", "w") as f:
            f.write(out)
        with open(f"{fp}/cloud-data", "w") as f:
            f.write("")


# testing only!!!!
if __name__ == "__main__":
    write(
        UserData(
            hostname="hostname",
            password="password",
            commands=["c", "o"],
            writes=[
                FileWrite("/etc/hostname", "root", "0644", "hostname\n"),
                FileWrite("/etc/hosts", "root", "0644", "hosts\n"),
            ],
        )
    )

import logging, tomllib, os
import mac
import vxml.switch, vxml.network
import ipaddress as ipa
import qcow2, ciso


def main(out="tmp/"):
    logger = logging.getLogger(__name__)
    with open("config.toml", "rb") as f:
        config = tomllib.load(f)
    devices = config["devices"]
    os.makedirs(out, exist_ok=True)

    name_type = {}
    name_conn = {}
    for dev in devices:
        name = dev["name"]
        type = dev["type"]
        name_type[name] = type

        name_conn[name] = dev["connections"]

    conn_pairs = []
    for name, type in name_type.items():
        if type == "switch":
            for peer in name_conn[name]:
                conn_pairs.append((name, peer))

    netws = {}
    for idx, (sw, cl) in enumerate(conn_pairs):
        idx += 1
        thoc = idx * 10
        bridge = f"virbr{thoc}"
        ip_addr = ipa.IPv4Interface(f"192.168.{thoc}.254/24")
        netw_name = f"{sw}-{cl}"
        with open(f"{out}/{netw_name}.xml", "w") as f:
            f.write(vxml.network.render(netw_name, bridge, ip_addr))
        netws[(sw, cl)] = (netw_name, ip_addr)

    mac_gen = mac.gen()
    for dev in devices:
        active_conn = []
        name = dev["name"]
        type = dev["type"]
        conns = dev["connections"]
        img_p = ""
        pkgs = []
        wfs = []
        cmds = []
        if type == "switch":
            img_p = os.path.abspath("images/switch")
            pkgs = ["arping"]
            cmds = [
                "sudo ip link set ens3 up",
                "sudo ip link set ens4 up",
                "sudo systemctl restart systemd-networkd",
                "sudo mkdir -p /etc/sysctl.d",
                'echo "net.ipv4.ip_forward=1" | sudo tee /etc/sysctl.d/99-ipforward.conf',
                "sudo sysctl --system",
            ]
            for idx, client in enumerate(conns):
                (netname, ip_inter) = netws[(name, client)]
                if netname is None:
                    logger.error(f"No network found for client {client}→ switch {name}")
                active_conn.append(netname)
                wfs.append(
                    [
                        f"/etc/systemd/network/{idx+1}0-ens{idx+3}.network",
                        f"[Match]\nName=ens{idx+3}\n\n[Network]\nAddress={ip_inter}",
                    ]
                )
        elif type == "client":
            img_p = os.path.abspath("images/dev")
            cmds = [
                "sudo ip link set ens3 up",
                "sudo systemctl restart systemd-networkd",
            ]
            for idx, conn in enumerate(conns):
                (netname, ip_inter) = netws[(conn, name)]
                if netname is None:
                    logger.error(f"No network found for client {name}→ switch {conn}")
                active_conn.append(netname)
                wfs.append(
                    [
                        f"/etc/systemd/network/{idx+1}0-ens{idx+2}.network",
                        f"[Match]\nName=ens{idx+2}\n\n[Network]\nAddress={ip_inter.network[2]}/24\nGateway={ip_inter.ip}",
                    ]
                )
        disk_p = qcow2.prep(logger, img_p, name)
        ciso_p = ciso.iso(
            logger,
            img_p,
            name,
            pkgs,
            wfs,
            cmds,
        )
        with open(f"{out}/{name}.xml", "w") as f:
            f.write(
                vxml.switch.render(
                    name,
                    disk_p,
                    ciso_p,
                    active_conn,
                    type,
                    mac_gen,
                )
            )


if __name__ == "__main__":
    main("tmp")

import ipaddress as ipa


def render(
    name: str,
    bridge: str,
    ipaddr: ipa.IPv4Interface,
):
    return f"""<network>
  <name>{name}</name>
  <forward mode='none'/>
  <bridge name='{bridge}' stp='on' delay='0'/>
  <ip address='{ipaddr.network[1]}' netmask='{ipaddr.netmask}'/>
</network>"""

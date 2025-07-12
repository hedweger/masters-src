import mac


def render(
    name: str,
    disk: str,
    seed: str | None,
    netw: list,
    type: str,
    mac_gen,
    nram=2048,
    vcpu=1,
):
    drives = f"""<disk type='file' device='disk'>
      <driver name='qemu' type='qcow2' cache='none'/>
      <source file='{disk}'/>
      <target dev='vda' bus='virtio'/>
    </disk>"""
    if seed:
        drives += f"""\n<disk type='file' device='cdrom'>
          <driver name='qemu' type='raw'/>
          <source file='{seed}'/>
          <target dev='hdc' bus='ide'/>
          <readonly/>
        </disk>"""
    ifaces = ""
    for net in netw:
        ifaces += f"""\n<interface type='network'>
          <mac address='{next(mac_gen)}'/>
          <source network='{net}'/>
          <model type='virtio'/>
        </interface>"""
    xml = f"""<domain type='kvm'>
  <name>{name}</name>
  <memory unit='MiB'>{nram}</memory>
  <vcpu placement='static'>{vcpu}</vcpu>
  <os>
    <type arch='x86_64' machine='pc'>hvm</type>
    <boot dev='hd'/>
  </os>
  <features>
    <acpi/>
    <apic/>
    <pae/>
  </features>
  <clock offset='utc'/>
  <on_poweroff>destroy</on_poweroff>
  <on_reboot>restart</on_reboot>
  <on_crash>destroy</on_crash>
  <devices>
    {drives}"""
    if type == "switch":
        xml += """<interface type='network'>
      <source network='default'/>
      <model type='virtio'/>
    </interface>"""
    xml += f"""{ifaces}
    <serial type='pty'>
      <target port='0'/>
    </serial>
    <console type='pty'>
      <target type='serial' port='0'/>
    </console>
  </devices>
</domain>
    """
    return xml

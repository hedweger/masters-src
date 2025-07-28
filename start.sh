#!/usr/bin/env bash
set -euo pipefail

declare -a NET_NAMES=( "sw1-pc1" "sw1-pc2" )
declare -a NET_XML=(   "tmp/netw-sw1-pc1.xml" "tmp/netw-sw1-pc2.xml" )

echo "==> Rebuilding networks..."

for i in "${!NET_NAMES[@]}"; do
  NAME="${NET_NAMES[i]}"
  XML="${NET_XML[i]}"

  echo "--- Processing network: $NAME (XML: $XML) ---"

  # 1) Destroy & undefine if they exist (ignore errors)
  echo "  -> destroying (if running): sudo virsh net-destroy $NAME"
  sudo virsh net-destroy "$NAME" 2>/dev/null || true

  echo "  -> undefining (if defined): sudo virsh net-undefine $NAME"
  sudo virsh net-undefine "$NAME" 2>/dev/null || true

  # 2) Define from XML
  echo "  -> defining from XML: sudo virsh net-define $XML"
  sudo virsh net-define "$XML"

  # 3) Start the network
  echo "  -> starting network: sudo virsh net-start $NAME"
  sudo virsh net-start "$NAME"

  # (optional) Enable autostart
  echo "  -> enabling autostart: sudo virsh net-autostart $NAME"
  sudo virsh net-autostart "$NAME"

  echo
done

declare -a VM_NAMES=( "pc1" "pc2" "sw1" )
declare -a VM_XML=(   "tmp/pc1/config.xml" "tmp/pc2/config.xml" "tmp/sw1/config.xml" )

echo "==> Rebuilding VMs..."

for i in "${!VM_NAMES[@]}"; do
  NAME="${VM_NAMES[i]}"
  XML="${VM_XML[i]}"

  echo "--- Processing VM: $NAME (XML: $XML) ---"

  # 1) Destroy (if running) & undefine (if defined)
  echo "  -> destroying domain (if running): sudo virsh destroy $NAME"
  sudo virsh destroy "$NAME" 2>/dev/null || true

  echo "  -> undefining domain (if defined): sudo virsh undefine $NAME"
  sudo virsh undefine "$NAME" 2>/dev/null || true

  # 2) Define from XML
  echo "  -> defining domain from XML: sudo virsh define $XML"
  sudo virsh define "$XML"

  # 3) Start the VM
  echo "  -> starting domain: sudo virsh start $NAME"
  sudo virsh start "$NAME"

  echo
done

echo "All networks and VMs have been redefined and started."

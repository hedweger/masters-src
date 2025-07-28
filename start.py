#!/usr/bin/env python3
import os
import subprocess

ROOT_DIR = "tmp"

def run_cmd(cmd):
    if cmd[0] == 'virsh':
        cmd = ['sudo'] + cmd
    try:
        output = subprocess.check_output(cmd, stderr=subprocess.STDOUT)
        return output.decode()
    except subprocess.CalledProcessError as e:
        return e.output.decode()

def process_xml(xml_file, command,vm_name):
    print(f"\nProcessing VM '{vm_name}' from {xml_file} ...")

    defined_vms = run_cmd(['virsh', command, '--all'])
    is_defined = vm_name in defined_vms

    running_vms = run_cmd(['virsh', command])
    is_running = vm_name in running_vms

    if is_running:
        print(f"Destroying running VM: {vm_name}")
        print(run_cmd(['virsh', 'destroy', vm_name]))

    if is_defined:
        print(f"Undefining VM: {vm_name}")
        print(run_cmd(['virsh', 'undefine', vm_name, '--remove-all-storage']))

    print(f"Defining VM from {xml_file}")
    print(run_cmd(['virsh', 'define', xml_file]))

    print(f"Starting VM: {vm_name}")
    print(run_cmd(['virsh', 'start', vm_name]))

def process_networks_first():
    networks_dir = os.path.join(ROOT_DIR, "networks")
    if not os.path.isdir(networks_dir):
        return
    for filename in sorted(os.listdir(networks_dir)):
        if filename.endswith(".xml"):
            xml_path = os.path.join(networks_dir, filename)
            vm_name = os.path.splitext(filename)[0]
            process_xml(xml_path, 'net-list', vm_name)

def process_devices():
    for entry in sorted(os.listdir(ROOT_DIR)):
        full_path = os.path.join(ROOT_DIR, entry)
        if os.path.isdir(full_path) and entry.startswith(('pc', 'sw')):
            for filename in sorted(os.listdir(full_path)):
                if filename.endswith('.xml'):
                    xml_path = os.path.join(full_path, filename)
                    vm_name = entry
                    process_xml(xml_path, 'list', vm_name)

def main():
    process_networks_first()
    process_devices()

if __name__ == "__main__":
    main()

# IED Network Simulator - Go Implementation

This directory contains the Go implementation of the IED Network Simulator, converted from Python.

## Overview

This tool manages virtual network devices and generates configurations for network simulation. It handles:

- Device management (RTUs and Switches)
- Network configuration generation
- Cloud-init configuration for VM initialization
- Libvirt XML generation for virtualization
- QCOW2 disk image management

## Usage

```bash
# Build the project
go build

# Run with config file (config.yaml should be in the same directory)
./ied-ns
```

## Architecture

The Go implementation maintains the same structure as the original Python version:

- **main.go**: Entry point
- **devman/**: Device management logic
- **device/**: Device models and types  
- **cinit/**: Cloud-init configuration generation
- **drive/**: Disk image handling
- **virtmac/**: MAC address generation

## Dependencies

- Go 1.21+
- `gopkg.in/yaml.v3` for YAML parsing
- System dependencies:
  - `genisoimage` for ISO generation
  - `cp` command for file operations

## Configuration

Uses the same `config.yaml` format as the Python version:

```yaml
network:
  address: "192.168.0.0/16"
  
rtus:
  - name: "pc1"
    address: "192.168.1.10"
  - name: "pc2"
    address: "192.168.1.20"

switches:
  - name: "sw1"
    address: ""
    connected:
      - to: "pc1"
      - to: "pc2"
```

## Templates

Go templates (in `templates/` directories) replace the original Jinja2 templates:

- `user-data.tmpl`: Cloud-init user data
- `network-config.tmpl`: Network configuration
- `virt_device.xml.tmpl`: Libvirt device XML
- `virt_network.xml.tmpl`: Libvirt network XML

## Output

The tool generates the same output structure as the Python version:

```
tmp/
├── pc1/
│   ├── debian-12-pc1.qcow2
│   ├── cloudinit.iso
│   ├── config.xml
│   └── seed/
│       ├── user-data
│       ├── network-config
│       └── meta-data
├── pc2/
│   └── ...
├── sw1-pc1.xml
└── sw1-pc2.xml
```
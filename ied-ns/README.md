# Go Implementation of Python main.py

This directory now contains a Go implementation that provides equivalent functionality to the existing Python `main.py` file.

## Files

### Go Implementation
- **main.go** - Main entry point, equivalent to Python main.py
- **types.go** - Go struct definitions for all data types
- **device_manager.go** - Core DeviceManager implementation
- **go.mod/go.sum** - Go module configuration

### Original Python Files
- **main.py** - Original Python implementation (preserved)
- **makedev/** - Python modules (preserved, used by Python version)

## Usage

### Run the Go version:
```bash
./main-go
```

### Run the Python version:
```bash
python3 main.py
```

## Functionality

Both implementations:
1. Read configuration from `../config.yaml`
2. Create RTU and Switch devices based on configuration
3. Generate network connections between devices
4. Create directory structure under `tmp/`
5. Generate libvirt XML configuration files

## Key Features Implemented in Go

- ✅ YAML configuration parsing
- ✅ Device management (RTUs and Switches)
- ✅ Network topology creation
- ✅ MAC address generation
- ✅ IP address management
- ✅ File and directory creation
- ✅ XML configuration generation

The Go implementation maintains the same interface and behavior as the Python version while providing better performance and static typing.
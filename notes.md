```
sudo ovs-vsctl add-br br0
sudo ip link set br0 up
sudo ovs-vsctl add-port br0 ens3
sudo ovs-vsctl add-port br0 ens4
sudo ip link set dev br0 up
```

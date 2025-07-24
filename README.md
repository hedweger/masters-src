- RTU
- [ ] sestavovani RTU bez .scd souboru (nachystat)
- [ ] muze but random hodnoty, nebo mirror z realnyho zarizeni
- Sitova vrstva
- [ ] pcap -> struktura
- [ ] musi jit vypnout aby to neslo poznat
- SCADA
- [ ] zatim chill, nemusim resit

- [ ] back kanal na ovladani switchu
- [ ] nejak vyresit nmap, resp. aby slo nacist ty sitovy struktury automaticky na real zarizeni
- [ ] nejakej analyzator na SCD file/struktury

1. parse config.toml file to get network layout + RTU units + etc.
2. prepare the correct number of VMs for the simulation:
    so far:
    - [X] one for each RTU (we could potentially run multiple RTUs in a single VM with concurrency)
    - [X] one for each switch (right now, maybe we could simulate networks inside a single VM, though I am not sure how this would work)
    - [ ] one for SCADA 
    - [ ] one for logs
    - [ ] one for WAN-GW (is this necessary? host device could serve this connection with correct firewall + NAT???)
3. slowly spin up each VM
    order of execution:
    - logging service
    - network set up (ideally, we also ping each SW to check network is correct)
    - RTU set up
    - SCADA connects to RTU servers
4. measure the network to see if everything works fine before letting the user access anything

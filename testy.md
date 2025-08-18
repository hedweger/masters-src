# Testovací plán 

## Společné sledované parametry (pro všechny testy)
- CPU:
  - Vzhledem k tomu, že síťový driver je kernel proces, tak nevím, jestli půjde sledovat, kolik přesně je přiřazeno CPU k síťovému procesu, ale dá se sledovat celkové využití.
  - Přes `virt-top` se dají sledovat stavy virtuálních strojů.
- RAM:
  - Využití celkem, případně swap (swap by VMky ideálně neměly potřebovat).
- Eth (Ens):
  - Propustnost
  - Latence (jitter?)
  - Hranice, kdy začne docházet k dropování paketů
  - Sledovat na switch zařízeních i na hostitelském systému
  
(Nástroj pro generování síťové zátěže: primárně iperf3, případně nějaký vlastní skript, který využívá iperf3 API)

## Test 1 – Maximální průchodnost switche při 1 až n paralelních spojeních 
- Zjistit, jak se mění dosažitelná propustnost a vytížení CPU se stoupajícím počtem připojených RTU.
- Parametry:
  - Počet spojení: 5, 10, 20, 30, 40, 50
  - Požadovaný bitrate: 1, 10, 20, 50, 100 MB/s
  - Velikost paketu: 256 B (TCP)

## Test 2 – Maximální velikost paketu při 1 až n paralelních spojeních
- Zjistit dopad velikosti paketu na propustnost a vytížení CPU.
- Parametry:
    - Počet spojení: 5, 10, 20, 30, 40, 50
    - Bitrate: 10 MB/s
    - Velikost paketu: 64 B, 128 B, 256 B, 512 B, 1024 B, 1500 B (TCP)

## Test 3 – Minimální konfigurace (vCPU, RAM) switche při 1 až n paralelních spojeních
- Zjistit minimální požadavky na vCPU a RAM pro dosažení stabilní propustnosti.
- Parametry:
    - Počet spojení: 5, 10, 20, 30, 40, 50
    - Požadovaný bitrate: 10 MB/s
    - Velikost paketu: 256 B (TCP)
    - Konfigurace VM: postupně budu snižovat vCPU a RAM virtuálního stroje switche, dokud nedojde k dropování paketů.

## Test 4 – Minimální konfigurace (vCPU, RAM) RTU při současném stress testu zařízení
- Zjistit minimální požadavky na VM parametry RTU zařízení, když zařízení dělá i něco jiného než jen přenosy.
- Parametry:
    - Počet spojení: 2 (1 SW + 2 RTU – sender a receiver)
    - Požadovaný bitrate: 10 MB/s
    - Velikost paketu: 256 B (TCP)
    - Konfigurace VM: postupně budu snižovat vCPU a RAM virtuálního stroje RTU, dokud nedojde k dropování paketů.
    - Současný stress test: například spuštění CPU zátěže na RTU (např. pomocí `stress` nebo `sysbench`).

## Test 5 – Minimální konfigurace (vCPU, RAM) switche a RTU při konfiguraci 50× RTU a 4 switche
- Parametry:
    - Počet spojení: 50 RTU + 4 SW (10 RTU na každý SW)
    - Požadovaný bitrate: 10 MB/s
    - Velikost paketu: 256 B (TCP)
- Nejdříve budu snižovat parametry RTU, dokud na RTU nedojde k dropování paketů na interfacech RTU, potom pro minimální konfiguraci RTU budu snižovat parametry switche, dokud opět nedojde k dropování paketů na interfacech switche.
- Cílem je zjistit, jak se mění propustnost a vytížení při vysokém počtu zařízení.

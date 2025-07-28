package device

import "fmt"

type MACGen struct {
	prefix  []byte
	current uint32
	max     uint32
}

func DefaultMACGenerator() *MACGen {
	return &MACGen{
		prefix:  []byte{0x52, 0x54, 0x00},
		current: 0,
		max:     0xFFFFFF,
	}
}

func (g *MACGen) Next() string {
	if g.current == g.max {
		panic("No more free MAC addresses available")
	}

	mac := make([]byte, 6)
	copy(mac, g.prefix)

	mac[3] = byte(g.current >> 16)
	mac[4] = byte(g.current >> 8)
	mac[5] = byte(g.current)
	g.current++

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

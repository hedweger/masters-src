package device

import (
	"fmt"
	"net"
	"strings"
)

type IPGen struct {
	network *net.IPNet
	prefix  string
	current uint32
	first   uint32
	last    uint32
}

func DefaultIPGenerator(cidr string) (*IPGen, error) {
	prefix := strings.Split(cidr, "/")[1]
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("Invalid CIDR: %w", err)
	}

	first, last := networkRange(network)
	first++
	last--

	return &IPGen{
		network: network,
		prefix:  prefix,
		current: first,
		first:   first,
		last:    last,
	}, nil
}

func (g *IPGen) Next() string {
	if g.current == g.last {
		panic(fmt.Sprintf("IP generator has reached the end of the range: %s", g.network.String()))
	}

	ip := make(net.IP, 4)
	ip[0] = byte(g.current >> 24)
	ip[1] = byte(g.current >> 16)
	ip[2] = byte(g.current >> 8)
	ip[3] = byte(g.current)

	g.current++
	return ip.String()
}

func (g *IPGen) NextWCidr() string {
	if g.current == g.last {
		panic(fmt.Sprintf("IP generator has reached the end of the range: %s", g.network.String()))
	}

	ip := make(net.IP, 4)
	ip[0] = byte(g.current >> 24)
	ip[1] = byte(g.current >> 16)
	ip[2] = byte(g.current >> 8)
	ip[3] = byte(g.current)

	g.current++
	return fmt.Sprintf("%s/%s", ip.String(), g.prefix)
}

func networkRange(network *net.IPNet) (uint32, uint32) {
	ip := network.IP.To4()
	if ip == nil {
		return 0, 0
	}

	first := uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
	mask := uint32(network.Mask[0])<<24 | uint32(network.Mask[1])<<16 | uint32(network.Mask[2])<<8 | uint32(network.Mask[3])
	last := first | ^mask

	return first, last
}

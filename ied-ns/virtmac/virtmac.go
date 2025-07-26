package virtmac

import "fmt"

// Generator generates MAC addresses under the prefix 52:54:00:12:xx:01,
// where xx bytes will increment with each call.
type Generator struct {
	current int
	wrapAt  int
}

// NewGenerator creates a new MAC address generator
func NewGenerator(start, wrapAt int) *Generator {
	if start < 0 {
		start = 0x22
	}
	if wrapAt <= 0 {
		wrapAt = 0x100
	}
	return &Generator{
		current: start,
		wrapAt:  wrapAt,
	}
}

// Next returns the next MAC address in the sequence
func (g *Generator) Next() string {
	mac := fmt.Sprintf("52:54:00:12:%02x:01", g.current)
	g.current = (g.current + 1) % g.wrapAt
	return mac
}

// Gen creates a generator with default values and returns a function
// that generates MAC addresses, similar to the Python version
func Gen() func() string {
	g := NewGenerator(0x22, 0x100)
	return g.Next
}
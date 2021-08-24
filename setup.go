package rebinding_protection

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/coredns/caddy"

	"sync"
)

var doOnce sync.Once

func init() { plugin.Register("rebinding_protection", setup) }

func setup(c *caddy.Controller) error {
	c.Next()

	block := New()

	if c.NextArg() {
		return plugin.Error("rebinding_protection", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		block.Next = next
		return block
	})

	return nil
}

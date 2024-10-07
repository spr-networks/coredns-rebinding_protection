package rebinding_protection

import (
	"context"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("rebinding_protection")

// Block is the block plugin.
type Block struct {
	Next plugin.Handler
	// our default block lists.

}

func New() *Block {
	return &Block{}
}

func (b *Block) ServeDNS(ctx context.Context, rw dns.ResponseWriter, r *dns.Msg) (int, error) {
	//this plugin is now deprecated, it is now part of 'block'
	return plugin.NextOrFailure(b.Name(), b.Next, ctx, rw, r)
}

// Name implements the Handler interface.
func (b *Block) Name() string { return "rebinding_protection" }

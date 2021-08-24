package rebinding_protection

import (
	"context"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
	"encoding/json"
	"strings"
	"net"
)

var log = clog.NewWithPlugin("rebinding_protection")

// Block is the block plugin.
type Block struct {
	Next plugin.Handler
	// our default block lists.

}

func New() *Block {
	return &Block{
	}
}

type EventData struct {
  Q []dns.Question
	A []dns.RR
}

type DNSEvent struct {
	dns.ResponseWriter
  data EventData
	delayedMsg *dns.Msg
}

func (i *DNSEvent) Write(b []byte) (int, error) {
	return i.ResponseWriter.Write(b)
}

func (i *DNSEvent) WriteMsg(m *dns.Msg) error {
	i.data.Q = m.Question
	i.data.A = m.Answer

	//delay the message until a decision has been made
	i.delayedMsg = m

	return nil
}

func (i *DNSEvent) String() string {
  x, _ := json.Marshal(i.data)
  return string(x)
}

//func (plugin JsonLog) ServeDNS(ctx context.Context, rw dns.ResponseWriter, r *dns.Msg) (c int, err error) {


type ResponseWriterDelay struct {
	dns.ResponseWriter
}


// ServeDNS implements the plugin.Handler interface.
func (b *Block) ServeDNS(ctx context.Context, rw dns.ResponseWriter, r *dns.Msg) (int, error) {

	event := &DNSEvent{
		ResponseWriter: rw,
	}

	c, err := b.Next.ServeDNS(ctx, event, r)

	for _, answer := range(event.data.A) {
		answerString := answer.String()
		parts := strings.Split(answerString, "\t")
		if len(parts) > 2 {
			rr_type := parts[len(parts)-2]
			if rr_type == "A" || rr_type == "AAAA" {
				ip := net.ParseIP(parts[len(parts)-1])

				// Need to block zero addresses as well
				_, zeroipv4, _ := net.ParseCIDR("0.0.0.0/32")
				_, zeroipv6, _ := net.ParseCIDR("::/32")

				// does golang handle ipv6 correctly?
				if ip.IsPrivate() || ip.IsLoopback() || ip.IsMulticast() ||
					ip.IsInterfaceLocalMulticast() || zeroipv4.Contains(ip) || zeroipv6.Contains(ip) {

					log.Infof("Blocking forward of %s, a local/private/multicast/zero IP address", ip.String())

					resp := new(dns.Msg)
					resp.SetRcode(r, dns.RcodeNameError)
					rw.WriteMsg(resp)

					return dns.RcodeNameError, nil
				}
			}

		}
	}

	event.ResponseWriter.WriteMsg(event.delayedMsg)

	return c, err
}

// Name implements the Handler interface.
func (b *Block) Name() string { return "rebinding_protection" }

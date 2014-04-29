package etcnet

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"net"
	"time"
)

type Client struct {
	*etcd.Client
}

// Time after which DNS entry expires from etcd (in seconds). Default is a day.
var DNS_ENTRY_TTL uint64 = 86400

// Creates a new client able to Dial using etcd as DNS store.
func NewClient(cluster []string) *Client {
	return &Client{etcd.NewClient(cluster)}
}

// Lookups host in etcd if not found lookups in DNS and stores the result in etcd.
func (c *Client) LookupHost(host string) (ips []string, err error) {
	// Lookup in etcd
	ips, err = c.Lookup(host)

	// Return if found in etcd
	if err == nil && len(ips) > 0 {
		return
	}

	// Otherwise lookup in DNS
	ips, err = net.LookupHost(host)
	if err != nil {
		return
	}

	// And save in etcd
	if err := c.register(host, ips, DNS_ENTRY_TTL); err != nil {
		log.Printf("Error saving addrs: %v", err)
	}
	return
}

// Lookups host in etcd.
func (c *Client) Lookup(host string) (addrs []string, err error) {
	// DNS host key-prefix
	prefix := c.hostprefix(host)

	// Get addresses from etcd
	res, err := c.Client.Get(prefix, false, false)
	if err != nil {
		return
	}
	plen := len(prefix) + 1 // length of prefix + /

	// Create a list of addrs from etcd response
	addrs = make([]string, res.Node.Nodes.Len())
	for n, r := range res.Node.Nodes {
		addrs[n] = r.Key[plen:]
	}
	return
}

// Set's a etcd cluster.
func (c *Client) SetCluster(cluster []string) bool {
	return c.Client.SetCluster(cluster)
}

// Dials address using etcdns Dialer
func (c *Client) Dial(network, address string) (net.Conn, error) {
	var d net.Dialer
	return c.internalDial(d, network, address)
}

// Dials address using etcdns Dialer with custom timeout
func (c *Client) DialTimeout(network, address string, timeout time.Duration) (conn net.Conn, err error) {
	d := net.Dialer{Timeout: timeout}
	return c.internalDial(d, network, address)
}

func (c *Client) internalDial(d net.Dialer, network, address string) (conn net.Conn, err error) {
	var e error
	var host, port string
	host, port, e = net.SplitHostPort(address)
	// Treat as "no port in address"
	if e != nil {
		host = address
	}
	// Lookup IP's
	addrs, err := c.LookupHost(host)
	if err != nil {
		return
	}

	// Range over IP's until connect
	for _, addr := range addrs {
		ip := addr
		// Set port from etcd if no any
		if port == "" {
			ip, port, err = net.SplitHostPort(addr)
			if err != nil {
				return
			}
		}
		conn, err = d.Dial(network, ip+":"+port)
		if err == nil {
			return
		}
		// we have an error
		// if eport is registered (port in etcd dns entry) - delete entry
		if err := c.Unregister(host, []string{addr}); err != nil {
			log.Printf("Error unregistering addrs: %v", err)
		}
	}

	return
}

// Writes given addresses in etcd under /dns/{host}/{ip} key.
// They will be possible to Dial like normal hosts.
func (c *Client) Register(host string, ips []string) error {
	return c.register(host, ips, 0)
}

func (c *Client) register(host string, ips []string, ttl uint64) (err error) {
	// Creates /dns/{host} key
	prefix := c.hostprefix(host)
	for _, ip := range ips {
		_, err = c.Client.Set(prefix+"/"+ip, "ok", ttl)
		if err != nil {
			return
		}
	}
	return
}

// Unregisters given IPs for host.
func (c *Client) Unregister(host string, ips []string) (err error) {
	// Creates /dns/{host} key
	prefix := c.hostprefix(host)
	for _, ip := range ips {
		_, err = c.Client.Delete(prefix+"/"+ip, true)
		if err != nil {
			return
		}
	}
	return
}

func (c *Client) hostprefix(host string) string {
	return "/dns/" + host
}

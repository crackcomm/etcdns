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

// Prefix of dns entries in etcd
var DNS_ETCD_PREFIX = "/dns/"

// Creates a new client able to Dial using etcd as DNS store.
func NewClient(cluster []string) *Client {
	return &Client{etcd.NewClient(cluster)}
}

// Looks for IP addresses (and ports if necessary) in etcd under /dns/{host} key.
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

// Looks for IP addresses in etcd using Lookup function.
// If nothing was found resolves from DNS using net.LookupHost and saves result in etcd.
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

// Set's a etcd cluster.
func (c *Client) SetCluster(cluster []string) bool {
	return c.Client.SetCluster(cluster)
}

// First lookups a IP address (and port if necessary) of host in etcd under /dns/{host} key.
// Then tries to connect.
func (c *Client) Dial(network, address string) (net.Conn, error) {
	var d net.Dialer
	return c.internalDial(d, network, address)
}

// Dials like Dial using custom timeout.
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

// Saves IP addresses (and ports if necessary) in etcd under /dns/{host}/{ip} key.
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

// Removes IP addresses (and ports if necessary) from etcd under /dns/{host}/{ip} keys.
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
	return DNS_ETCD_PREFIX + host
}

package etcnet

import (
	"net"
	"time"
)

var DefaultClient = NewClient([]string{"127.0.0.1:4001"})

// Looks for IP addresses (and ports if necessary) in etcd under /dns/{host} key.
func Lookup(host string) (ips []string, err error) {
	return DefaultClient.Lookup(host)
}

// Looks for IP addresses in etcd using Lookup function.
// If nothing was found resolves from DNS using net.LookupHost and saves result in etcd.
func LookupHost(host string) (ips []string, err error) {
	return DefaultClient.LookupHost(host)
}

// Set's a etcd cluster.
func SetCluster(cluster []string) bool {
	return DefaultClient.SetCluster(cluster)
}

// First lookups a IP address (and port if necessary) of host in etcd under /dns/{host} key.
// Then tries to connect.
func Dial(network, address string) (net.Conn, error) {
	return DefaultClient.Dial(network, address)
}

// Dials like Dial using custom timeout.
func DialTimeout(network, address string, timeout time.Duration) (conn net.Conn, err error) {
	return DefaultClient.DialTimeout(network, address, timeout)
}

// Saves IP addresses (and ports if necessary) in etcd under /dns/{host}/{ip} key.
func Register(host string, ips []string) error {
	return DefaultClient.Register(host, ips)
}

// Removes IP addresses (and ports if necessary) from etcd under /dns/{host}/{ip} keys.
func Unregister(host string, ips []string) error {
	return DefaultClient.Unregister(host, ips)
}

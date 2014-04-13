package etcnet

import (
	"net"
	"time"
)

var DefaultClient = NewClient([]string{"127.0.0.1:4001"})

// Lookups host in etcd if not found lookups in DNS and stores the result in etcd.
func LookupHost(host string) (ips []string, err error) {
	return DefaultClient.LookupHost(host)
}

// Set's a etcd cluster.
func SetCluster(cluster []string) bool {
	return DefaultClient.SetCluster(cluster)
}

// Dials address using etcdns Dialer
func Dial(network, address string) (net.Conn, error) {
	return DefaultClient.Dial(network, address)
}

// Dials address using etcdns Dialer with custom timeout
func DialTimeout(network, address string, timeout time.Duration) (conn net.Conn, err error) {
	return DefaultClient.DialTimeout(network, address, timeout)
}

// Writes given addresses in etcd under /dns/{host}/{ip} key.
// They will be possible to Dial like normal hosts.
func Register(host string, ips []string) error {
	return DefaultClient.Register(host, ips)
}

// Unregisters given IPs for host.
func Unregister(host string, ips []string) error {
	return DefaultClient.Unregister(host, ips)
}

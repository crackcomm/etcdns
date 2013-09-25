package etcnet

import (
	"log"
	"net"
)

type endpoint struct {
	ip   string
	ping int64 // milliseconds

}

var DNS_EXPIRE uint64 = 60 * 60 * 24 // in seconds

func LookupHost(host string) ([]string, error) {
	return lookupHost(host)
}

// Lookups host in etcd if not any found resolves with DNS and saves in etcd
func lookupHost(host string) (addrs []string, err error) {
	addrs, err = etcdResolve(host)
	if err == nil && len(addrs) > 0 {
		return
	}
	addrs, err = dnsResolve(host)
	if err != nil {
		return
	}
	if serr := etcdSaveAddrs(host, addrs); serr != nil {
		log.Printf("Error saving addrs: %v", serr)
	}
	return
}

// Saves resolved addrs of host to etcd
func etcdSaveAddrs(host string, addrs []string) error {
	if Debug {
		t := ben.Start("etcd.save")
		defer t.End()
	}
	prefix := "/dns/" + host + "/"
	for _, addr := range addrs {
		if _, err := Client.Set(prefix+addr, "ok", DNS_EXPIRE); err != nil {
			return err
		}
	}
	return nil
}

// Resolves host addresses using etcd
func etcdResolve(host string) (list []string, err error) {
	if Debug {
		t := ben.Start("etcd.resolve")
		defer t.End()
	}
	// Get addresses from etcd
	prefix := "/dns/" + host + "/"
	res, err := Client.Get(prefix)
	if err != nil || len(res) == 0 {
		return
	}
	plen := len(prefix)
	// Create a list of addrs
	list = make([]string, len(res))
	for n, r := range res {
		list[n] = r.Key[plen:]
	}
	return
}

// Resolves host using DNS
func dnsResolve(host string) ([]string, error) {
	if Debug {
		t := ben.Start("dns.resolve")
		defer t.End()
	}
	return net.LookupHost(host)
}

package etcnet

import (
	"net"
	"time"
)

func Dial(network, address string) (net.Conn, error) {
	var d net.Dialer
	return internalDial(d, network, address)
}

func DialTimeout(network, address string, timeout time.Duration) (conn net.Conn, err error) {
	d := net.Dialer{Timeout: timeout}
	return internalDial(d, network, address)
}

func internalDial(d net.Dialer, network, address string) (conn net.Conn, err error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return
	}
	addrs, err := LookupHost(host)
	if err != nil {
		return
	}
	for _, addr := range addrs {
		conn, err = d.Dial(network, addr+":"+port)
		if err == nil {
			return
		}
	}
	return
}

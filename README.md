etcdns
======

Extended Go `net` dialing with DNS caching backed by [etcd](https://github.com/coreos/etcd) a highly-available key value store.

example
=======

Service discovery using etcdns:

```Go

import (
	net "github.com/crackcomm/etcdns/net"
)

// on server & client
func init() {
	net.SetCluster([]string{"etcd.com:4001"})
}

func server() {
	go net.Register("secret.service", []string{"10.0.0.1:5051"}) // register service in etcd
	http.ListenAndServe("10.0.0.1:5051", nil)                 // start listening for connections
}

func client() {
	conn, err := net.Dial("tcp", "secret.service")
	if err != nil {
		fmt.Printf("Error dialing secret.service")
	}
	defer conn.Close()
	// we are now connected to 10.0.0.1:5051
}
```

godoc
=====

[godoc](http://godoc.org/github.com/crackcomm/etcdns/net)


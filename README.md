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
	go http.ListenAndServe("10.0.0.1:5051", nil)              // some listener
	net.Register("secret.service", []string{"10.0.0.1:5051"}) // register service in etcd
}

func client() {
	conn, err := net.Dial("tcp", "secret.service")            // will dial to 10.0.0.1:5051
	if err != nil {
		fmt.Printf("Error dialing secret.service")
	}
	defer conn.Close()
	// Use connection
	// ...
}
```

godoc
=====

[godoc](http://godoc.org/github.com/crackcomm/etcdns/net)


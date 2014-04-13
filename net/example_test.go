package etcnet

import (
	"fmt"
	"net/http"
)

func ExampleDial() {
	// Dial using etcdns
	conn, err := Dial("tcp", "golang.org:80")
	if err != nil {
		fmt.Printf("Error dialing golang.org")
	}
	defer conn.Close()

	// Write some HTTP request
	_, err = conn.Write([]byte("GET / HTTP/1.1\nHost: golang.org\n"))
	if err != nil {
		fmt.Printf("Error writing request to golang.org")
	}

	// now we can read request...
}

func ExampleDial_serviceDiscovery() {
	// on the server
	SetCluster([]string{"public.etcd.ip:4001"})
	go http.ListenAndServe("my.public.ip:5051", nil)          // some listener
	Register("secret.service", []string{"my.public.ip:5051"}) // register service in etcd

	// on the client
	SetCluster([]string{"public.etcd.ip:4001"})
	conn, err := Dial("tcp", "secret.service")
	if err != nil {
		fmt.Printf("Error dialing golang.org")
	}
	defer conn.Close()
	// Use connection
	// ...
}

func ExampleSetCluster() {
	SetCluster([]string{"127.0.0.1:4001", "127.0.0.1:4002", "127.0.0.1:4003", "127.0.0.1:4004"})
}

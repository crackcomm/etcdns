package etcnet

import (
	"log"
	"net/http"
	"time"
)

func ExampleDial() {
	conn, err := Dial("tcp", "secret.service")
	if err != nil {
		log.Printf("Error dialing golang.org")
	}
	defer conn.Close()

	// Now we are connected to 10.0.0.1:5051
	// we can write a HTTP request...
	_, err = conn.Write([]byte("GET / HTTP/1.1\nHost: secret.service\n"))
}

func ExampleDialTimeout() {
	// Resolves a `secret.service` to IP and port from etcd
	// Aborts if it's not possible under 10 seconds
	conn, err := DialTimeout("tcp", "secret.service", 10*time.Second)
	if err != nil {
		log.Printf("Error dialing golang.org")
	}
	defer conn.Close()

	// Now we are connected to 10.0.0.1:5051
	// we can write a HTTP request...
	_, err = conn.Write([]byte("GET / HTTP/1.1\nHost: secret.service\n"))
}

func ExampleLookup() {
	ips, err := Lookup("secret.service")
	if err != nil {
		log.Fatalf("lookup error:", err)
	}
	log.Println("ips =>", ips)
	// ips => []string{"10.0.0.1:5051"}
}

func ExampleLookupHost() {
	ips, err := LookupHost("secret.service")
	if err != nil {
		log.Fatalf("lookup error:", err)
	}
	log.Println("ips =>", ips)
	// ips => []string{"10.0.0.1:5051"}
}

func ExampleRegister() {
	go Register("secret.service", []string{"10.0.0.1:5051"}) // register service in etcd
	http.ListenAndServe("10.0.0.1:5051", nil)                // start listening for connections
}

func ExampleUnregister() {
	// when registering a service in etcd
	// defer unregister on clean exit
	defer func() {
		Unregister("secret.service", []string{"10.0.0.1:5051"})
	}()
}

func ExampleSetCluster() {
	SetCluster([]string{"127.0.0.1:4001", "127.0.0.1:4002", "127.0.0.1:4003", "127.0.0.1:4004"})
}

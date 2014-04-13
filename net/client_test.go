package etcnet

import (
	"os"
	"strings"
	"testing"
)

var googleIps = []string{"173.194.43.38", "173.194.43.39", "173.194.43.40", "173.194.43.41"}

func TestSetCluster(t *testing.T) {
	cluster := strings.Split(os.Getenv("ETCD_CLUSTER"), ",")
	if ok := SetCluster(cluster); !ok {
		t.Errorf("Error syncing cluster %s (set env ETCD_CLUSTER)", cluster)
	}
}

func TestRegister(t *testing.T) {
	// google.com
	if err := Register("google", googleIps); err != nil {
		t.Error(err)
	}
}

func TestLookupHost(t *testing.T) {
	ips, err := LookupHost("google")
	if err != nil {
		t.Error(err)
		return
	}
	if len(ips) != len(googleIps) {
		t.Errorf("There should be exactly %d IP addresses got %d", len(googleIps), len(ips))
		return
	}

	// Check if all IP's were registered
	for _, gip := range googleIps {
		found := false
		for _, ip := range ips {
			if !found && ip == gip {
				found = true
			}
		}
		if !found {
			t.Errorf("IP %s is not present in etcd", gip)
		}
	}
}

func TestUnregister(t *testing.T) {
	err := Unregister("google", []string{googleIps[0]})
	if err != nil {
		t.Error(err)
		return
	}
	ips, err := LookupHost("google")
	if err != nil {
		t.Error(err)
		return
	}
	if len(ips) < len(googleIps)-1 {
		t.Errorf("There should be exactly %d IP addresses got %d", len(googleIps)-1, len(ips))
		return
	}

	// Check if all IP's were registered
	for _, ip := range ips {
		if ip == googleIps[0] {
			t.Errorf("IP %s was not unregistered", ip)
		}
	}
}

func TestDial(t *testing.T) {
	conn, err := Dial("tcp", "google:80")
	if err != nil {
		t.Errorf("Dial error: %v", err)
		return
	}
	defer conn.Close()
}

// tests removing not available dns entry
func TestDial_remove(t *testing.T) {
	// no.available
	err := Register("no.available", []string{"127.0.0.1:9889"})
	if err != nil {
		t.Error(err)
		return
	}

	ips, err := LookupHost("no.available")
	if err != nil {
		t.Error(err)
		return
	}
	if len(ips) != 1 || ips[0] != "127.0.0.1:9889" {
		t.Errorf("no.available host was not registered")
		return
	}

	// it wont be successful
	_, err = Dial("tcp", "no.available")
	if err == nil {
		t.Errorf("Connection to 127.0.0.1:9889 was successful - it can't be!")
		return
	}

	ips, _ = LookupHost("no.available")
	if len(ips) != 0 {
		t.Errorf("no.available entry was not removed")
		return
	}
}

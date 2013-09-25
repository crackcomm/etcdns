package etcnet

import "github.com/coreos/go-etcd/etcd"

var Client = etcd.NewClient()

func SetCluster(cluster string) bool {
	return Client.SetCluster(cluster)
}

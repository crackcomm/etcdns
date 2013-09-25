package etcnet

import (
	"github.com/crackcomm/utils/benchmark"
	"github.com/coreos/go-etcd/etcd"
)

var (
	Debug = false
	ben = bench.New("etcnet")
)

func OpenDebug() {
	Debug = true	
	etcd.OpenDebug()
}

func CloseDebug() {
	Debug = false	
	etcd.CloseDebug()
}

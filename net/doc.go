/*

Extended net dialing using https://github.com/coreos/etcd a highly-available key value store.

Difference between etcnet.Lookup and etcnet.LookupHost is:

1. etcnet.Lookup just lookups in etcd

2. etcnet.LookupHost lookups in etcd using etcnet.Lookup, if nothing was found lookups in DNS and saves result in etcd

*/
package etcnet

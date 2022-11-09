package etcd

import (
	"fmt"
	"testing"
)

var ETCD Etcd

func init() {
	etcd, err := NewEtcd(`{"addr":"172.16.102.109:2379,172.16.102.109:2380","password":"","time-out":5}`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ErrEmptyAddr)
	ETCD = etcd

	_, err = ETCD.Put("/etcd_test/1", "test1")
	if err != nil {
		return
	}

	go func(ETCD Etcd) {
		res := ETCD.WatchPrefix("/etcd_test/")
		for i := range res {
			for _, event := range i.Events {
				fmt.Printf("type:%s key:%s value: %s\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}(ETCD)

	go func(ETCD Etcd) {
		res := ETCD.Watch("/etcd_test/1")
		for i := range res {
			for _, event := range i.Events {
				fmt.Printf("type:%s key:%s value: %s\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}(ETCD)
}

func TestClient_Put(t *testing.T) {
	_, err := ETCD.Put("/etcd_test/2", "test")
	if err != nil {
		return
	}
}

func TestClient_Get(t *testing.T) {
	res, err := ETCD.Get("/etcd_test/1")
	if err != nil {
		fmt.Println(err)
	}

	for _, i2 := range res.Kvs {
		fmt.Println(string(i2.Value))
	}
}

func TestClient_GetPrefix(t *testing.T) {
	res, err := ETCD.GetPrefix("/etcd_test/")
	if err != nil {
		fmt.Println(err)
	}

	for _, i2 := range res.Kvs {
		fmt.Println(string(i2.Value))
	}
}

func TestClient_DelGet(t *testing.T) {
	_, err := ETCD.DelGet("/etcd_test/1")
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_PutLease(t *testing.T) {
	lease, _ := ETCD.LeaseGrant(5)
	opt := ETCD.WithLease(lease.ID)
	_, err := ETCD.PutLease("/etcd_test/lease", "lease_test", opt)
	if err != nil {
		fmt.Println(err)
	}
}

package rpc

import etcdClientV3 "go.etcd.io/etcd/client/v3"

type Rpc interface {
	Client() *Grpc
	Register(*etcdClientV3.Client) (*ServiceRegister, error)
	Run()
}

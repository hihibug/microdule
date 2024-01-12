package rpc

import (
	grpcs "github.com/hihibug/microdule/rpc/grpc"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Rpc interface {
	Client() any
	Register(*etcdClientV3.Client) (*grpcs.ServiceRegister, error)
	Run() error
	Close()
}

type Grpc struct {
	RpcSrv       *grpc.Server
	Config       *Config
	EtcdRegister *grpcs.ServiceRegister
}

package rpc

import (
	"fmt"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
)

type Grpc struct {
	RpcSrv       *grpc.Server
	Config       *Config
	EtcdRegister *ServiceRegister
}

func NewGrpc(c *Config, opt ...grpc.ServerOption) Rpc {
	grpcServer := grpc.NewServer(opt...)

	return &Grpc{
		RpcSrv: grpcServer,
		Config: c,
	}
}

func (g *Grpc) Client() *Grpc {
	return g
}

func (g *Grpc) Register(etcd *etcdClientV3.Client) (*ServiceRegister, error) {
	return NewRpcServiceRegister(etcd,
		g.Config.Name,
		fmt.Sprint(g.Config.IP, ":", g.Config.Addr),
		5,
	)
}

func (g *Grpc) Run() {

	address := fmt.Sprintf(":%d", g.Config.Addr)

	lis, err := net.Listen("tcp", address)

	if err != nil {
		panic("grpc net listen error :" + err.Error())
	}

	err = g.RpcSrv.Serve(lis)

	if err != nil {
		panic(err)
	}
}

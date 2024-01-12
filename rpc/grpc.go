package rpc

import (
	"fmt"
	"net"

	grpcs "github.com/hihibug/microdule/rpc/grpc"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func NewGrpc(c *Config, opt ...grpc.ServerOption) Rpc {
	grpcServer := grpc.NewServer(opt...)

	return &Grpc{
		RpcSrv: grpcServer,
		Config: c,
	}
}

func (g *Grpc) Client() any {
	return g
}

func (g *Grpc) Register(etcd *etcdClientV3.Client) (*grpcs.ServiceRegister, error) {
	return grpcs.NewRpcServiceRegister(etcd,
		g.Config.Name,
		fmt.Sprint(g.Config.IP, ":", g.Config.Addr),
		5,
	)
}

func (g *Grpc) Run() error {

	address := fmt.Sprintf(":%d", g.Config.Addr)

	lis, err := net.Listen("tcp", address)

	defer g.EtcdRegister.Close()

	if err != nil {
		panic("grpc net listen error :" + err.Error())
	}

	go g.EtcdRegister.ListenLeaseRespChan()

	err = g.RpcSrv.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

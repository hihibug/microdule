package rpc

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Grpc struct {
	RpcSrv *grpc.Server
	Config Config
}

func NewGrpc(c Config) Rpc {
	grpcServer := grpc.NewServer()

	return &Grpc{
		RpcSrv: grpcServer,
		Config: c,
	}
}

func (g *Grpc) GetGrpc() *Grpc {
	return g
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

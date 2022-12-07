package rpc

type Rpc interface {
	GetGrpc() *Grpc
	Run()
}

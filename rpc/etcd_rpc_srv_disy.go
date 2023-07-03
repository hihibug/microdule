package rpc

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"sync"
)

type EtcdResolverBuilder struct {
	etcdClient *clientv3.Client
}

type etcdResolver struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cc         resolver.ClientConn
	etcdClient *clientv3.Client
	scheme     string
	ipPool     sync.Map
}

// NewRpcServiceDiscovery  新建发现服务
func NewRpcServiceDiscovery(prefix string, etcdCli *clientv3.Client) (*grpc.ClientConn, error) {

	etcdResolverBuilder := NewEtcdResolverBuilder(etcdCli)
	//resolver.Register(etcdResolverBuilder)
	conn, err := grpc.Dial(
		"etcd:///"+prefix,
		grpc.WithResolvers(etcdResolverBuilder),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewEtcdResolverBuilder(etcdCli *clientv3.Client) *EtcdResolverBuilder {
	return &EtcdResolverBuilder{
		etcdClient: etcdCli,
	}
}

func (erb *EtcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	// 获取 etcd 中服务保存的ip列表
	res, err := erb.etcdClient.Get(context.Background(), target.URL.Path, clientv3.WithPrefix())

	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	es := &etcdResolver{
		cc:         cc,
		etcdClient: erb.etcdClient,
		ctx:        ctx,
		cancel:     cancelFunc,
		scheme:     target.URL.Path,
	}

	// 将获取到的ip和port保存到本地的map
	for _, kv := range res.Kvs {
		es.store(kv.Key, kv.Value)
	}

	// 更新拨号里的ip列表
	es.updateState()

	// 监听etcd中的服务是否变化
	go es.watcher()

	return es, nil
}

func (erb *EtcdResolverBuilder) Scheme() string {
	return "etcd"
}

func (e *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {
	//global.GO_LOG.Info("etcd 发现器正在更新")
}

func (e *etcdResolver) Close() {
	//global.GO_LOG.Info("etcd 发现器关闭链接")
	e.cancel()
}

func (e *etcdResolver) watcher() {
	watchChan := e.etcdClient.Watch(context.Background(), e.scheme, clientv3.WithPrefix())

	for {
		select {
		case val := <-watchChan:
			for _, event := range val.Events {
				switch event.Type {
				case 0: // 0 是有数据增加
					e.store(event.Kv.Key, event.Kv.Value)
					//global.GO_LOG.Info(e.scheme + "服务数量增加")
					e.updateState()
				case 1: // 1是有数据减少
					e.del(event.Kv.Key)
					//global.GO_LOG.Info(e.scheme + "服务数量减少")
					e.updateState()
				}
			}
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *etcdResolver) store(k, v []byte) {
	e.ipPool.Store(string(k), string(v))
}

func (s *etcdResolver) del(key []byte) {
	s.ipPool.Delete(string(key))
}

func (e *etcdResolver) updateState() {
	var addrList resolver.State
	e.ipPool.Range(func(k, v interface{}) bool {
		tA, ok := v.(string)
		if !ok {
			return false
		}
		addrList.Addresses = append(addrList.Addresses, resolver.Address{Addr: tA})
		return true
	})

	e.cc.UpdateState(addrList)
}

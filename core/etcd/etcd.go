package etcd

import (
	"context"
	client3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type (
	// Client etcd 链接
	Client struct {
		Cli  *client3.Client
		Time time.Duration
	}

	Etcd interface {
		Get(name string) (res *client3.GetResponse, err error)                                   // Get 取值
		GetPrefix(name string) (res *client3.GetResponse, err error)                             // GetPrefix 前缀取值
		Put(name, value string) (res *client3.PutResponse, err error)                            // Put 存值
		PutLease(name, value string, opt client3.OpOption) (res *client3.PutResponse, err error) // PutLease 租约存值
		DelGet(name string) (res *client3.DeleteResponse, err error)                             // DelGet 删值
		Watch(name string) (watchCh client3.WatchChan)                                           // Watch 监听
		WatchPrefix(name string) (watchCh client3.WatchChan)                                     // WatchPrefix 前缀监听
		LeaseGrant(i int) (leaseResp *client3.LeaseGrantResponse, err error)                     // LeaseGrant 租约
		WithLease(id client3.LeaseID) client3.OpOption                                           // WithLease 租约
	}
)

func (c *Client) Get(name string) (res *client3.GetResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Get(ctx, name)
	ConFunc()
	return
}

func (c *Client) GetPrefix(name string) (res *client3.GetResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Get(ctx, name, client3.WithPrefix())
	ConFunc()
	return
}

func (c *Client) Put(name, value string) (res *client3.PutResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Put(ctx, name, value)
	ConFunc()
	return
}

func (c *Client) PutLease(name, value string, opt client3.OpOption) (res *client3.PutResponse, err error) {
	res, err = c.Cli.Put(context.Background(), name, value, opt)
	return
}

func (c *Client) DelGet(name string) (res *client3.DeleteResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Delete(ctx, name)
	ConFunc()
	return
}

func (c *Client) Watch(name string) (watchCh client3.WatchChan) {
	watchCh = c.Cli.Watch(context.Background(), name)
	return
}

func (c *Client) WatchPrefix(name string) (watchCh client3.WatchChan) {
	watchCh = c.Cli.Watch(context.Background(), name, client3.WithPrefix())
	return
}

func (c *Client) LeaseGrant(i int) (leaseResp *client3.LeaseGrantResponse, err error) {
	leaseResp, err = c.Cli.Grant(context.Background(), int64(i+1))
	return
}

func (c *Client) WithLease(id client3.LeaseID) client3.OpOption {
	return client3.WithLease(id)
}

// NewEtcd 创建etcd链接
func NewEtcd(conf string) (Etcd, error) {

	configs, err := DeConfig(conf)

	if err != nil {
		return nil, err
	}

	addr := strings.Split(configs.Addr, ",")

	cli, err := client3.New(client3.Config{
		Endpoints:   addr,
		Password:    configs.Password,
		DialTimeout: time.Duration(configs.TimeOut) * time.Second,
	})

	if err != nil {
		return nil, err
	}

	ctx, ConFunc := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, "test-ping", "go-ping-etcd-test")
	ConFunc()

	if err != nil {
		return nil, err
	}

	return &Client{
		Cli:  cli,
		Time: time.Duration(configs.TimeOut) * time.Second,
	}, nil
}

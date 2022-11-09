package etcd

import (
	"context"
	client3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

// Client etcd 链接
type Client struct {
	Cli  *client3.Client
	Time time.Duration
}

// NewEtcd 创建etcd链接
func NewEtcd(conf string) (*Client, error) {

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

// Get 取值
func (c *Client) Get(name string) (res *client3.GetResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Get(ctx, name)
	ConFunc()
	return
}

// GetPrefix 前缀取值
func (c *Client) GetPrefix(name string) (res *client3.GetResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Get(ctx, name, client3.WithPrefix())
	ConFunc()
	return
}

// Put 存值
func (c *Client) Put(name, value string) (res *client3.PutResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Put(ctx, name, value)
	ConFunc()
	return
}

// PutLease 租约存值
func (c *Client) PutLease(name, value string, opt client3.OpOption) (res *client3.PutResponse, err error) {
	res, err = c.Cli.Put(context.Background(), name, value, opt)
	return
}

// DelGet 删值
func (c *Client) DelGet(name string) (res *client3.DeleteResponse, err error) {
	ctx, ConFunc := context.WithTimeout(context.Background(), c.Time)
	res, err = c.Cli.Delete(ctx, name)
	ConFunc()
	return
}

// Watch 监听
func (c *Client) Watch(name string) (watchCh client3.WatchChan) {
	watchCh = c.Cli.Watch(context.Background(), name)
	return
}

// WatchPrefix 前缀监听
func (c *Client) WatchPrefix(name string) (watchCh client3.WatchChan) {
	watchCh = c.Cli.Watch(context.Background(), name, client3.WithPrefix())
	return
}

// LeaseGrant 租约
func (c *Client) LeaseGrant(i int) (leaseResp *client3.LeaseGrantResponse, err error) {
	leaseResp, err = c.Cli.Grant(context.Background(), int64(i+1))
	return
}

// WithLease 租约
func (c *Client) WithLease(id client3.LeaseID) client3.OpOption {
	return client3.WithLease(id)
}

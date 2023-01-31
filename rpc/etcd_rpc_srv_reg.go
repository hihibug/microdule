package rpc

import (
	"context"
	"log"
	"os"
	"strconv"

	etcdClientV3 "go.etcd.io/etcd/client/v3"
)

//ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	cli           *etcdClientV3.Client                        //etcd client
	leaseID       etcdClientV3.LeaseID                        //租约ID
	keepAliveChan <-chan *etcdClientV3.LeaseKeepAliveResponse //租约keepalieve相应chan
	key           string                                      //key
	val           string                                      //value
}

// NewRpcServiceRegister 新建注册服务
func NewRpcServiceRegister(etcd *etcdClientV3.Client, key, val string, lease int64) (*ServiceRegister, error) {
	ser := &ServiceRegister{
		cli: etcd,
		key: key,
		val: val,
	}

	//申请租约设置时间keepalive并注册服务
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}

	return ser, nil
}

//设置租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	//设置租约时间
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	key := "/" + s.key + "/" + strconv.Itoa(int(resp.ID))
	_, err = s.cli.Put(context.Background(), key, s.val, etcdClientV3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan
	log.Printf("grpc  name: %s", s.key)
	log.Printf("lease id  : %d", s.leaseID)
	return nil
}

//ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		_ = leaseKeepResp
		// log.Println("续约成功", leaseKeepResp)
	}
	log.Println("lease close")
	os.Exit(0)
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Println("grpc close :", s.leaseID)
	return s.cli.Close()
}

package modular

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

func (etcd *Client) CallRegister(host string, appName string, value []byte, ttl int) error {
	ctx := context.Background()
	prefix := etcd.keyPrefix(host, appName)
	ttlResp, err := etcd.Backend.Grant(context.TODO(), int64(ttl))
	if err != nil {
		etcd.logger.Errorf("etcd: register client.Grant(%v) error(%v)", ttl, err)
		return err
	} else {
		_, err = etcd.Backend.Put(ctx, prefix, string(value), clientv3.WithLease(ttlResp.ID))
		if err != nil {
			etcd.logger.Errorf("etcd: register client.Put(%v) key(%s) error(%v)", prefix, host, err)
			return err
		}
	}
	return nil
}

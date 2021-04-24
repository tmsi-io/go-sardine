package modular

import (
	"context"
	"encoding/json"
	"github.com/tmsi-io/go-sardine/naming/datadefine"
	"go.etcd.io/etcd/clientv3"
)

// 查询带状态接口
func (etcd *Client) DiscoveryService(appName string) []datadefine.AppInstance {
	var ins []datadefine.AppInstance
	prefix := etcd.keyPrefix("", appName)
	ctx := context.Background()
	resp, err := etcd.Backend.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		etcd.logger.Errorf("Etcd get error:[%v]", err)
		return ins
	}
	for _, ev := range resp.Kvs {
		var serv datadefine.BaseAppInfo
		err := json.Unmarshal(ev.Value, &serv)
		if err != nil {
			etcd.logger.Errorf("Unmarshal failed:[%v]", err)
			continue
		} else {
			ins = append(ins, datadefine.AppInstance{BaseInfo: serv})
		}
	}
	return ins
}

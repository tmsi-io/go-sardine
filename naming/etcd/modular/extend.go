package modular

import (
	"context"
	"encoding/json"
	"github.com/tmsi-io/go-sardine/naming/datadefine"
)

// GetExtendInfo get the apps extend config
func (etcd *Client) GetExtendInfo(key string, appName string) (result datadefine.ExSetting, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.timeOut)
	defer cancel()
	keyPrefix := etcd.exKeyPrefix(key, appName)
	if response, err := etcd.Backend.Get(ctx, keyPrefix); err != nil {
		return
	} else {
		if len(response.Kvs) == 1 && string(response.Kvs[0].Key) == keyPrefix {
			if err := json.Unmarshal(response.Kvs[0].Value, &result); err != nil {
			} else {
				ok = true
				return
			}
		}
	}
	return
}

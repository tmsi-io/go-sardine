package modular

import (
	"fmt"
)

// keyPrefix Get BaseApp info keyPrefix
func (etcd *Client) keyPrefix(key string, appName string) string {
	return fmt.Sprintf("/%s/%s/%s", DiscoveryPrefix, appName, key)
}

// exKeyPrefix get extend key prefix
func (etcd *Client) exKeyPrefix(key string, appName string) string {
	return fmt.Sprintf("/%s/%s/%s", ExtendPrefix, appName, key)
}

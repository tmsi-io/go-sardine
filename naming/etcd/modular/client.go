package modular

import (
	"go.etcd.io/etcd/clientv3"
	"net/url"
	"strings"
	"time"
)

func GetClient() *Client {
	return &cl
}

func (etcd *Client) Init(_url string) (err error) {
	pUrl, err := url.Parse(_url)
	if err != nil {
		return err
	}
	password, _ := pUrl.User.Password()
	etcd.Backend, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(pUrl.Host, ","),
		DialTimeout: 5 * time.Second,
		Username:    pUrl.User.Username(),
		Password:    password,
	})
	return err
}

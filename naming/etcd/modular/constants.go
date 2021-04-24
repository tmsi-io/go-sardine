package modular

import (
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Client struct {
	Backend *clientv3.Client
	logger  *logrus.Entry
	timeOut time.Duration
}

var cl Client

const (
	DiscoveryPrefix = "Discovery"
	ExtendPrefix    = "Extend"
)

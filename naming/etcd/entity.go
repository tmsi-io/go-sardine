package etcd

import (
	"github.com/tmsi-io/go-sardine/naming/datadefine"
	"github.com/tmsi-io/go-sardine/naming/etcd/modular"
)

type Etcd struct {
	appName string
	server  string
}

func (e *Etcd) SetEnv(app string, addr string) {
	e.appName = app
	e.server = addr
}
func (e *Etcd) GetClient() error {
	cl := modular.GetClient()
	return cl.Init(e.server)
}

func (e *Etcd) GetBaseInfoKey() string {
	panic("implement me")
}

func (e *Etcd) GetExtendInfoKey() string {
	panic("implement me")
}

func (e *Etcd) Register() {
	panic("implement me")
}

func (e *Etcd) Discovery(apps []string) {
	panic("implement me")
}

func (e *Etcd) GetOneApp(choice []func()) (datadefine.AppInstance, error) {
	panic("implement me")
}

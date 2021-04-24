package modular

import "go-sardine/naming/datadefine"

type ZooKeeper struct{}

func (z *ZooKeeper) SetEnv(app string, url string) {
	panic("implement me")
}

func (z *ZooKeeper) GetClient() error {
	panic("implement me")
}

func (z *ZooKeeper) GetBaseInfoKey() string {
	panic("implement me")
}

func (z *ZooKeeper) GetExtendInfoKey() string {
	panic("implement me")
}

func (z *ZooKeeper) Register() {
	panic("implement me")
}

func (z *ZooKeeper) Discovery(apps []string) {
	panic("implement me")
}

func (z *ZooKeeper) GetOneApp(choice []func()) (datadefine.AppInstance, error) {
	panic("implement me")
}

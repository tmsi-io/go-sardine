package consul

import "go-sardine/naming/datadefine"

type Consul struct{}

func (c *Consul) SetEnv(app string, addr string) {

}

func (c *Consul) GetClient(addr string) error {
	panic("implement me")
}

func (c *Consul) GetBaseInfoKey() string {
	panic("implement me")
}

func (c *Consul) GetExtendInfoKey() string {
	panic("implement me")
}

func (c *Consul) Register() {
	panic("implement me")
}

func (c *Consul) Discovery(apps []string) {
	panic("implement me")
}

func (c *Consul) GetOneApp(choice []func()) (datadefine.AppInstance, error) {
	panic("implement me")
}

package eureka

import (
	"errors"
	"fmt"
	"github.com/tmsi-io/go-sardine/naming/datadefine"
	"github.com/tmsi-io/go-sardine/naming/eureka/modular"
	"net/url"
	"strings"
)

// eureka
type eureka struct {
	appName string
	server  string
}

func (e *eureka) SetEnv(app string, url string) {
	e.appName = app
	e.server = url
}

func (e *eureka) GetClient() error {
	_url := strings.Split(e.server, ",")
	if len(_url) == 0 {
		return errors.New("there was no eureka server config. ")
	} else {
		for _, addr := range _url {
			if _, err := url.Parse(addr); err != nil {
				return fmt.Errorf("there was config error for '%s', %v. \n", addr, err)
			}
		}
		return modular.InitEureka(_url)
	}
}

func (e *eureka) GetBaseInfoKey() string {
	panic("implement me")
}

func (e *eureka) GetExtendInfoKey() string {
	panic("implement me")
}

func (e *eureka) Register() {
	panic("implement me")
}

func (e *eureka) Discovery(apps []string) {
	panic("implement me")
}

func (e *eureka) GetOneApp(choice []func()) (datadefine.AppInstance, error) {
	panic("implement me")
}

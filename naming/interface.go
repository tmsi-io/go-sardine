package naming

import "go-sardine/naming/datadefine"

type Interface interface {
	SetEnv(app string, url string)
	// get register client
	GetClient() error
	// get base info key config
	GetBaseInfoKey() string
	// get extend info key config
	GetExtendInfoKey() string
	// do register for self
	Register()
	// do discovery for target
	Discovery(apps []string)
	// get one app by choice
	GetOneApp(choice []func()) (datadefine.AppInstance, error)
}

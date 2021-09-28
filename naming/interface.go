package naming

import "go-sardine/naming/datadefine"

type Interface interface {
	SetEnv(app string, url string)
	// GetClient get register client
	GetClient() error
	// GetBaseInfoKey get base info key config
	GetBaseInfoKey() string
	// GetExtendInfoKey get extend info key config
	GetExtendInfoKey() string
	// Register do register for self
	Register()
	// Discovery do discovery for target
	Discovery(apps []string)
	// GetOneApp get one app by choice
	GetOneApp(choice []func()) (datadefine.AppInstance, error)
}

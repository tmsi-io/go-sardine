package modular

import "fmt"

const GetUrlsMaxReTry = 5

const EurekaStatusUP = "UP"
const EurekaStatusDown = "DOWN"

type EurekaAppInfo struct {
	hostname        string
	HttpPort        int
	ServicePriority int
	ServiceVersion  string
}

func (eInfo *EurekaAppInfo) String() string {
	return fmt.Sprintf("%s:%d:%d:%s",
		eInfo.hostname, eInfo.HttpPort, eInfo.ServicePriority,
		eInfo.ServiceVersion,
	)
}

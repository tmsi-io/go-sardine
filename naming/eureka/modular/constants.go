package modular

import "fmt"

const GetUrlsMaxReTry = 5

const EurekaStatusUP = "UP"
const EurekaStatusDown = "DOWN"

type EurekaAppInfo struct {
	hostname string
	HttpPort int
	Priority int
	Version  string
}

func (eInfo *EurekaAppInfo) String() string {
	return fmt.Sprintf("%s:%d:%d:%s",
		eInfo.hostname, eInfo.HttpPort, eInfo.Priority,
		eInfo.Version,
	)
}

package modular

import (
	"strconv"
	"time"
)

const (
	EurekaHeartBeatTime     = time.Second * 70 // 心跳保活时长
	EurekaHeartBeatInterval = time.Second * 60 // 心跳间隔时长
	EurekaHealthCheckURL    = "/sidecar/health"
	EurekaStatusCheckURL    = "/sidecar/status"
	LocalHTTPScheme         = "http://"
)

// 初始化并保活
// RegisterAndHeartBeat
func RegisterAndHeartBeat(app, hostname, clientIP string, cPort int, SrvUP bool, instanceID string) error {
	var strStatus string
	if SrvUP {
		strStatus = EurekaStatusUP
	} else {
		strStatus = EurekaStatusDown
	}
	ins := Instance{
		InstanceId:       instanceID,
		HostName:         hostname,
		App:              app,
		Port:             &Port{Port: cPort, Enable: "true"},
		IPAddr:           clientIP, // sometimes ip better, like in docker
		VipAddress:       clientIP,
		SecureVipAddress: clientIP,
		HealthCheckUrl:   LocalHTTPScheme + clientIP + ":" + strconv.Itoa(cPort) + EurekaHealthCheckURL, // {"status": "UP|DOWN"}
		StatusPageUrl:    LocalHTTPScheme + clientIP + ":" + strconv.Itoa(cPort) + EurekaStatusCheckURL, // 状态端口
		HomePageUrl:      LocalHTTPScheme + clientIP + ":" + strconv.Itoa(cPort),
		Status:           strStatus,
		DataCenterInfo: &DataCenterInfo{
			Name:  "MyOwn",
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
	}
	err := eureka.RegisterInstance(&ins)
	if err != nil {
		return err
	}
	go func() {
		eureka.SendHeartBeat(&ins, EurekaHeartBeatTime)
	}()
	return nil
}

// StatusOutOfService
func StatusOutOfService(appID string, instanceID string) error {
	return eureka.OutOfServiceInstance(appID, instanceID)
}

// StatusUP
func StatusUP(appID string, instanceID string) error {
	return eureka.UPInstance(appID, instanceID)
}

// StatusDown
func StatusDown(appID string, instanceID string) error {
	return eureka.DownInstance(appID, instanceID)
}

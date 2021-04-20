package modular

import "net/http"

type ServerInfo struct {
	Urls       []string       //服务器列表
	DataCenter DataCenterInfo //"class: com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo"
}

type InstanceInfo struct {
	Id          string     //实例ID
	HostName    string     //
	AppName     string     //
	Ip          string     //
	BizPort     int        //
	StatusPort  int        //
	HealthCheck string     //
	StatusCheck string     //
	HomePageUrl string     //
	Status      string     //
	server      ServerInfo //
	eureka      *Eureka
}

type Eureka struct {
	ServiceUrls []string
	Client      *http.Client
	Json        bool
}

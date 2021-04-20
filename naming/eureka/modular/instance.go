package modular

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (ins *InstanceInfo) Init() {
	var err error
	ins.eureka, err = NewEureka(ins.server.Urls, nil)
	if err != nil {
		log.Fatalln("Initial Eureka obj Failed.")
	}
}

// NewInstance
func (ins *InstanceInfo) NewInstance() (Instance, error) {
	if ins.BizPort == 0 || ins.Ip == "" || ins.AppName == "" || ins.Id == "" || ins.StatusPort == 0 {
		return Instance{}, errors.New("loss argument. ")
	}
	if ins.HealthCheck == "" {
		ins.HealthCheck = "http://" + ins.Ip + ":" + strconv.Itoa(ins.StatusPort) + "/health"
	}
	if ins.StatusCheck == "" {
		ins.StatusCheck = "http://" + ins.Ip + ":" + strconv.Itoa(ins.StatusPort) + "/status"
	}
	if ins.HomePageUrl == "" {
		ins.HomePageUrl = "http://" + ins.Ip + ":" + strconv.Itoa(ins.StatusPort)
	}
	if ins.Status == "" {
		ins.Status = "DOWN"
	}
	_ins := Instance{
		InstanceId:       ins.Id,
		HostName:         ins.HostName,
		App:              ins.AppName,
		Port:             &Port{Port: ins.BizPort, Enable: "true"},
		IPAddr:           ins.Ip,
		VipAddress:       ins.Ip,
		SecureVipAddress: ins.Ip,
		HealthCheckUrl:   ins.HealthCheck,
		StatusPageUrl:    ins.StatusCheck,
		HomePageUrl:      ins.HomePageUrl,
		Status:           ins.Status,
		DataCenterInfo: &DataCenterInfo{
			Name:  ins.server.DataCenter.Name,
			Class: ins.server.DataCenter.Class,
		},
	}
	return _ins, nil
}

// 注册app
// RegisterApp
func (ins *InstanceInfo) RegisterApp(_ins Instance) bool {
	err := ins.eureka.RegisterInstance(&_ins)
	if err != nil {
		log.Fatalln(fmt.Sprintf("GenerateInstance Error: %v. ", err))
		return false
	}
	return true
}

// 服务器端默认超时时间是30秒
// SendHeartBeat
func (ins *InstanceInfo) SendHeartBeat(_ins Instance, interval int) {
	ins.eureka.SendHeartBeat(&_ins, time.Second*20)
}

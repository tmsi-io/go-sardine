package modular

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (e *Eureka) pickServerUrl() string {
	/*
		检查服务端可用性
	*/
	urls := e.ServiceUrls
	l := len(urls)
	if l == 0 {
		return ""
	}
	if l == 1 {
		if checkIp(urls[0] + "/apps") {
			return urls[0]
		}
	}
	r := rand.Intn(l)
	url := urls[r]
	if checkIp(url + "/apps") {
		return url
	} else {
		for i := 0; i < l; i++ {
			if checkIp(urls[i] + "/apps") {
				return urls[i]
			}
		}
	}
	return ""
}

// 注册实例
func (e *Eureka) RegisterInstance(i *Instance) error {
	/*
		POST   /eureka/apps/appID
		Body:  JSON/XML payload
		HTTPCode: 204 on success
	*/
	url := e.pickServerUrl()
	if len(url) == 0 {
		return errors.New("Url Too Short. ")
	}
	i.Init()
	// Instance数据构建
	app := App{
		Instance: i,
	}
	data, err := json.Marshal(&app)
	if err != nil {
		return err
	}
	if req, err := http.NewRequest("POST", url+"/apps/"+i.App, bytes.NewReader(data)); err != nil {
		return err
	} else {
		req.Header.Set("Content-Type", "application/json")
		res, err := e.Client.Do(req)
		if err != nil {
			return err
		}
		defer func() { _ = res.Body.Close() }()
		if res.StatusCode != 204 {
			return errors.New(strconv.Itoa(res.StatusCode))
		}
		return nil
	}
}

// 发送心跳
func (e *Eureka) SendHeartBeat(i *Instance, duration time.Duration) {
	go func() {
		urls := e.ServiceUrls
		if len(urls) == 0 {
			fmt.Println("missing eureka url.")
		}
		for {
			e.HearBeat(urls, i)
			time.Sleep(duration)
		}
	}()
}

func (e *Eureka) HearBeat(urls []string, i *Instance) {
	req, err := http.NewRequest("PUT", urls[0]+"/apps/"+i.App+"/"+i.InstanceId, nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := e.Client.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer func() { _ = res.Body.Close() }()
		statusCode := res.StatusCode
		if statusCode != 200 {
			fmt.Printf("Eureka HearBeat return %d \n ", statusCode)
			_ = e.RegisterInstance(i)
		}
	}
}

func (e *Eureka) GetApp(appID string) (*Application, error) {
	/*
		GET /eureka/apps/appID
	*/
	urls := e.ServiceUrls
	l := len(urls)
	url := urls[rand.Intn(l)]
	req, err := http.NewRequest("GET", url+"/apps/"+appID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	result := Application{}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response StatusCode is not 200，but " + strconv.Itoa(resp.StatusCode))
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &result, nil
}

// 获取APP url列表
func (e *Eureka) GetAppUrls(appID string) []string {
	/*
		通过实例ID查询服务url列表
	*/
	var maxTry int
	app, err := e.GetApp(appID)
	for err != nil {
		fmt.Println(err)
		if maxTry >= GetUrlsMaxReTry {
			return []string{}
		}
		app, err = e.GetApp(appID)
		maxTry++
	}
	var urls []string
	for _, ins := range app.Application.Instance {
		if ins.Status == "UP" {
			url := ins.IPAddr + ":" + strconv.Itoa(ins.Port.Port)
			if checkIp(ins.HealthCheckUrl) {
				urls = append(urls, url)
			}
		}
	}
	return urls
}

// 删除实例
func (e *Eureka) DelInstance(appID, instanceID string) error {
	/*
		DELETE /eureka/apps/appID/instanceID
	*/
	urls := e.ServiceUrls
	if len(urls) == 0 {
		return errors.New("missing eureka url. ")
	}
	req, err := http.NewRequest("DELETE", urls[0]+"/apps/"+appID+"/"+instanceID, nil)
	if err != nil {
		return err
	}
	res, err := e.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode == 200 {
		return nil
	} else {
		return errors.New("unknown error. ")
	}
}

// 将某个实例设置为下线
func (e *Eureka) OutOfServiceInstance(appID, instanceID string) error {
	/*
		PUT /eureka/apps/appID/instanceID/status?value=OUT_OF_SERVICE
	*/

	urls := e.ServiceUrls
	if len(urls) == 0 {
		return errors.New("missing eureka url. ")
	}
	req, err := http.NewRequest("PUT", urls[0]+"/apps/"+appID+"/"+instanceID+"/status?value=OUT_OF_SERVICE", nil)
	if err != nil {
		return err
	}
	res, err := e.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode == 200 {
		return nil
	} else {
		return errors.New("unknown error. ")
	}
}

func (e *Eureka) DownInstance(appID, instanceID string) error {
	/*
		PUT /eureka/apps/appID/instanceID/status?value=Down
	*/

	urls := e.ServiceUrls
	if len(urls) == 0 {
		return errors.New("missing eureka url. ")
	}
	req, err := http.NewRequest("PUT", urls[0]+"/apps/"+appID+"/"+instanceID+"/status?value=DOWN", nil)
	if err != nil {
		return err
	}
	res, err := e.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode == 200 {
		return nil
	} else {
		return errors.New("unknown error. ")
	}
}

func (e *Eureka) UPInstance(appID, instanceID string) error {
	/*
		PUT /eureka/apps/appID/instanceID/status?value=UP
	*/

	urls := e.ServiceUrls
	if len(urls) == 0 {
		return errors.New("missing eureka url. ")
	}
	req, err := http.NewRequest("PUT", urls[0]+"/apps/"+appID+"/"+instanceID+"/status?value=UP", nil)
	if err != nil {
		return err
	}
	res, err := e.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode == 200 {
		return nil
	} else {
		return errors.New("unknown error. ")
	}
}

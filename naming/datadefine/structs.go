package datadefine

import (
	"context"
	"math"
	"sync"
)

/*
	注册信息内容
*/
type BaseAppInfo struct {
	Name        string   `json:"Name"`        // 应用名称, 所属应用集群
	ID          string   `json:"AppID"`       // 应用ID， 不能重复,格式： AppName.Env.host.md5hash, 消费MQ或者kafka时，用作consumer ID
	CreateTime  int64    `json:"CreateTime"`  // Redis记录创建时间
	RefreshTIme int64    `json:"RefreshTIme"` // 心跳时间
	Zone        string   `json:"Zone"`        // 机房
	Env         string   `json:"Env"`         // 环境   Dev ; Fat ; Uat ; Pro
	Color       string   `json:"Color"`
	HostName    string   `json:"HostName"`     // 主机名
	AddrS       []string `json:"AddrS"`        // 地址列表
	Version     string   `json:"Version"`      // 版本号
	Weight      int      `json:"Weight"`       // 权重，程序默认
	CurrLoad    int      `json:"CurrLoad"`     // 当前负载数量
	MaxLoad     int      `json:"MAXTaskCount"` // 最大负载数量
	Metadata    []byte   `json:"Metadata"`     // 元数据列表
}

/*
	扩展控制部分
*/
/*
Key同注册服务字段，前缀改为Ex
例如：Ex. RTSPRedirect. hDiENNxem ； Ex. RTSPRedirect. node5:9556
Value字段
*/
type ExSetting struct {
	Weight    int               `json:"Weight"`       //权重，扩展。若有扩展，采用扩展字段。
	Status    Status            `json:"Status"`       //状态    0：不可用 ； 1可用
	MaxLoad   int               `json:"MAXTaskCount"` // 最大负载数
	MetaData  []byte            `json:"MetaData"`
	IDCStatus map[string]Status `json:"IDCStatus"` // 新加字段,当前服务想对于各个机房的使能状态
	IDCWeight map[string]int    `json:"IDCWeight"` // 新加字段,当前服务想对于各机房的权重和使用优先级
	IsExist   bool              `json:"-"`
}

type AppServiceList struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	AppName      string
	Services     map[string]*AppInstance
	ServicesLock sync.RWMutex
	Status       int32
}

type AppInstance struct {
	CacheKey   string
	BaseInfo   BaseAppInfo
	Extend     ExSetting
	PickWeight map[string]int
}

func (pThis AppInstance) GetPickWeight(zone string) int {
	if weight, ok := pThis.PickWeight[zone]; ok {
		return weight
	}
	return 0 // 默认未设置相关机房设置时,设置为1 or 0?
}

// 获取权重
// GetWeight
func (pThis AppInstance) GetWeight(zone string) int {
	if weight, ok := pThis.Extend.IDCWeight[zone]; ok {
		return weight
	}
	return pThis.Extend.Weight
}

func (pThis AppInstance) GetIDCWeight() map[string]int {
	if pThis.Extend.IsExist {
		return pThis.Extend.IDCWeight
	} else {
		return nil
	}
}

func (pThis AppInstance) GetMaxLoad() int {
	if pThis.Extend.IsExist {
		return pThis.Extend.MaxLoad
	} else {
		return pThis.BaseInfo.MaxLoad
	}
}

func (pThis AppInstance) Available() bool {
	if pThis.Extend.IsExist {
		return pThis.Extend.Status == StatusEnable
	}
	return true
}

func (pThis AppInstance) AvailableByZone(zone string) bool {
	if pThis.Extend.IsExist {
		if status, ok := pThis.Extend.IDCStatus[zone]; ok {
			return status == StatusEnable
		}
		return pThis.Extend.Status == StatusEnable
	}
	return true
}

// 刷新服务对于已经设置的各机房的负载
func (pThis AppInstance) RefreshZoneWeight() {
	var IDCWeight map[string]int
	if pThis.PickWeight == nil {
		pThis.PickWeight = make(map[string]int)
	}
	if pThis.Extend.IsExist && pThis.Extend.IDCWeight != nil {
		IDCWeight = pThis.Extend.IDCWeight
	} else {
		IDCWeight = map[string]int{pThis.BaseInfo.Zone: pThis.Extend.Weight}
	}
	for zone, weight := range IDCWeight {
		var pickWeight int
		var gift int
		// 当服务所在机房与设置机房一致时,加上一层, 确保负载设置后,流量尽量压到指定机房,
		// pickWeight最大范围==0~Sqrt(maxload), 按设置的limit超过范围需要200^2=40000最大负载服务才会有坏处理
		// 设置时,第一梯度机房(边缘侧)权重设置为0-30?, 后续梯度机房,根据实际设置
		if pThis.BaseInfo.Zone == zone {
			gift = 1
		}
		// 负载已满
		if pThis.GetMaxLoad()-pThis.BaseInfo.CurrLoad < 0 {
			pickWeight = 0
		} else {
			//(maxlimit-优先级)＊math.Sqrt(MaxLoad-LocalCount)
			pickWeight = (WeightLimit - weight + 1 + gift*WeightLimit) * int(math.Sqrt(float64(pThis.GetMaxLoad()-pThis.BaseInfo.CurrLoad)))

		}
		pThis.PickWeight[zone] = pickWeight
	}
}

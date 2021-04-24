package datadefine

type NamingType string

/*
	服务发现类型选择
*/
const (
	Redis  NamingType = "Redis"
	Etcd   NamingType = "Etcd"
	Eureka NamingType = "Eureka"
	Consul NamingType = "Consul" // 暂不实现
)

/*
	服务发现与注册时间配置
*/
const (
	RegisterTTL   = 15 // 注册ttl时间,单位秒
	ServiceExpire = RegisterTTL
)

/*
	服务状态
*/

type Status int32

const (
	StatusDisable     Status = 0 // 服务禁用
	StatusEnable      Status = 1 // 服务启用
	ServiceStatusStop Status = -1
)

var (
	// WeightLimit was the Weight limit config, uint config, smaller to higher
	WeightLimit int = 100
	// Priority was current service weight
	Weight int = WeightLimit / 2
)

type Config struct {
	Main   string `json:"Main"`   // 主Naming
	Redis  string `json:"Redis"`  // Redis 访问地址
	Etcd   string `json:"Etcd"`   // Etcd 访问地址
	Eureka string `json:"Eureka"` // Eureka 访问地址
}

package datadefine

type NamingType string

const (
	Redis  NamingType = "Redis"
	Etcd   NamingType = "Etcd"
	Eureka NamingType = "Eureka"
	Consul NamingType = "Consul"
)

const (
	RegisterTTL   = 15
	ServiceExpire = RegisterTTL
)

/*
	服务状态
*/

type Status int32

const (
	StatusDisable     Status = 0 //
	StatusEnable      Status = 1 //
	ServiceStatusStop Status = -1
)

var (
	// WeightLimit was the weight limit config, uint config, smaller to higher
	WeightLimit int = 100
)

type Config struct {
	Main   string `json:"Main"`   //
	Redis  string `json:"Redis"`  // Redis
	Etcd   string `json:"Etcd"`   // Etcd
	Eureka string `json:"Eureka"` // Eureka
}

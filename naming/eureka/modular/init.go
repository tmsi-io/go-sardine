package modular

import "log"

var eureka *Eureka

// 初始化
// InitEureka
func InitEureka(serverUrls []string) (err error) {
	eureka, err = NewEureka(serverUrls, nil)
	if err != nil {
		log.Print(err)
		return
	}
	return
}

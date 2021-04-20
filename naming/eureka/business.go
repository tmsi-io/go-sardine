package eureka

// NewInstance
// 初始化新的 eureka instance
//func Register(appName string, bizPort int, serverUrl string, statusPort int) error {
//
//	if err := DoRegister(instanceID, hostname, clientIP); err != nil {
//		return err
//	} else {
//		return nil
//	}
//}
//
//func DoRegister(instanceID string, hostname string, clientIP string) (err error) {
//	if err = modular.InitEureka(strings.Split(EurekaURL, ",")); err != nil {
//		return err
//	}
//	return modular.RegisterAndHeartBeat(EurekaAppName, hostname, clientIP, EurekaPort, true, instanceID)
//}

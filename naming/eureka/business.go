package eureka

//func GeneralInstanceID(app string, bizPort int, statusPort int) string{
//
//}
//
//// Register
//// eureka instance
//func Register(appName string, bizPort int, serverUrl string, statusPort int) error {
//	instanceID:=GeneralInstanceID(appName, bizPort, statusPort)
//	host:=utils.GetHostName()
//	if err := DoRegister(instanceID, host, clientIP, bizPort); err != nil {
//		return err
//	} else {
//		return nil
//	}
//}
//
//func DoRegister(instanceID string, hostname string, clientIP string, ) (err error) {
//	if err = modular.InitEureka(strings.Split(EurekaURL, ",")); err != nil {
//		return err
//	}
//	return modular.RegisterAndHeartBeat(EurekaAppName, hostname,clientIP,  , true, instanceID)
//}

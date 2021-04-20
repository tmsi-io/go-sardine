package modular

// 按app名称发现服务
// DiscoveryByAppID
func DiscoveryByAppID(appID string) []string {
	if eureka != nil {
		return eureka.GetAppUrls(appID)
	}
	return []string{}
}

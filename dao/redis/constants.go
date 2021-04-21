package redis

// error content
const (
	ErrorConnLoss   = "Connection loss. "
	ErrorConfigHost = "Can't get host and port in config. "
	ErrorConfigDB   = "Can't Parse DB Config. "
)

// config fields
const (
	ConfigMaxRetries = "MaxRetries"
	ConfigPoolSize   = "PoolSize"
)

// Networks, second
const (
	CfgDialTimeout = 10
)

// Pools , second
const (
	CfgPoolTimeout = 1
	CfgIdleTimeout = 20
)

// Sentinel constants
const (
	DefaultSentinel = "mymaster"
	SentinelAddr1   = "redis-sentinel"
	SentinelAddr2   = "redis+sentinel"
)

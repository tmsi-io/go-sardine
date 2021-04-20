module go-sardine

go 1.15

require (
	github.com/go-redis/redis/v7 v7.4.0
	github.com/prometheus/client_golang v1.10.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tmsi-io/go-sardine v0.0.0-20210420024515-c2f9fc8f6cc3
)

replace github.com/tmsi-io/go-sardine v0.0.0-20210420024515-c2f9fc8f6cc3 => ./

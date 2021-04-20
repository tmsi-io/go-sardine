package modular

import (
	"errors"
	"net/http"
)

func NewEureka(serverUrls []string, client *http.Client) (*Eureka, error) {
	if len(serverUrls) == 0 {
		return nil, errors.New("missing eureka url. ")
	}
	if client == nil {
		client = http.DefaultClient
	}
	eureka := &Eureka{
		ServiceUrls: serverUrls,
		Client:      client,
	}
	return eureka, nil
}

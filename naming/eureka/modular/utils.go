package modular

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"
)

var Timeout = 500 * time.Millisecond
var roundTrip = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	DisableKeepAlives: true,
	Proxy:             http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          300,
	MaxIdleConnsPerHost:   50,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var Client = &http.Client{Timeout: Timeout, Transport: roundTrip}

// checkIp
// check ip was ok
func checkIp(url string) bool {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	resp, err := Client.Do(req)
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

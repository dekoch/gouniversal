package clientinfo

import (
	"net"
	"net/http"
)

type ClientInfo struct {
	UserAgent string
	IP        string
}

func Get(r *http.Request) ClientInfo {

	var c ClientInfo
	c.UserAgent = r.UserAgent()

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	c.IP = ip

	return c
}

func String(r *http.Request) string {

	c := Get(r)

	return "UA: " + c.UserAgent + ", IP: " + c.IP
}

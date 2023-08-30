package main

import (
	"errors"
	"net/http/httputil"
	"net/url"
)

type Port string

type HostMap struct {
	host  map[string]Port
	admin string
	re    map[string]*httputil.ReverseProxy
}

func NewHostMap() *HostMap {
	return &HostMap{
		host:  make(map[string]Port),
		re:    make(map[string]*httputil.ReverseProxy),
		admin: "localhost:8000",
	}
}

func (hm *HostMap) GetHost(port Port) (response string, err error) {
	for host, p := range hm.host {
		if p == port {
			return host, nil
		}
	}
	return "", errors.New("host not found")
}

func (hm *HostMap) GetHostsArray() (response []string) {
	response = append(response, hm.admin)
	for host := range hm.host {
		response = append(response, host)
	}
	return
}

func (hm *HostMap) GetPort(host string) (response Port, err error) {
	if p, ok := hm.host[host]; ok {
		return p, nil
	}
	return "", errors.New("port not found")
}

func (hm *HostMap) GetProxy(host string) (response *httputil.ReverseProxy, err error) {
	if p, ok := hm.re[host]; ok {
		return p, nil
	}
	return nil, errors.New("proxy not found")
}

func (hm *HostMap) Set(host string, port Port) {
	if port == "admin" {
		hm.admin = host
		return
	}
	hm.host[host] = port
	hm.re[host] = httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "localhost:" + string(port)})
}

func (hm *HostMap) SetAll(hosts map[string]Port) {
	for host, port := range hosts {
		hm.Set(host, port)
	}
}

func (hm *HostMap) Delete(host string) {
	delete(hm.host, host)
	delete(hm.re, host)
}

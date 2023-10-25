package main

import (
	"encoding/json"
	"fmt"
)

type Prune struct {
	model    string
	ip       string
	port     string
	user     string
	password string

	httpRequest *HttpRequest
	path        string
}

type Client struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
}

func NewPrune(model, ip, port, user, password string) *Prune {
	var path string
	if model == ModelController {
		path = "v2/api/site/default/clients"
	} else {
		path = "proxy/network/v2/api/site/default/clients"
	}
	if port == "" {
		port = "443"
	}
	prune := &Prune{
		model:       model,
		ip:          ip,
		port:        port,
		user:        user,
		password:    password,
		httpRequest: NewHttpRequest(10),
		path:        path,
	}
	return prune
}

func (p *Prune) Run() error {
	err := p.Login()
	if err != nil {
		return err
	}
	fmt.Println("login success")
	macs, err := p.GetOfflineClients()
	if err != nil {
		return err
	}
	for _, mac := range macs {
		fmt.Println(mac)
	}
	err = p.Logout()
	if err != nil {
		return err
	}
	fmt.Println("logout success")
	return nil
}

func (p *Prune) Login() error {
	params := map[string]interface{}{
		"username": p.user,
		"password": p.password,
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	url := fmt.Sprintf("https://%s:%s/api/auth/login", p.ip, p.port)

	_, err := p.httpRequest.Request(url, "POST", params, headers)
	if err != nil {
		return err
	}
	p.httpRequest.GetCookies()
	return nil
}

func (p *Prune) Logout() error {
	url := fmt.Sprintf("https://%s:%s/api/auth/logout", p.ip, p.port)
	headers := make(map[string]string)
	token, ok := p.httpRequest.headers["X-Csrf-Token"]
	if ok && len(token) > 0 {
		headers["X-Csrf-Token"] = token[0]
	}

	_, err := p.httpRequest.Request(url, "POST", nil, headers)
	if err != nil {
		return err
	}
	return nil
}

func (p *Prune) GetOfflineClients() ([]string, error) {
	params := "onlyNonBlocked=true&includeUnifiDevices=true&withinHours=0"
	url := fmt.Sprintf("https://%s:%s/%s/history?%s", p.ip, p.port, p.path, params)
	rspStr, err := p.httpRequest.Request(url, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	var clients []*Client
	err = json.Unmarshal(rspStr, &clients)
	if err != nil {
		return nil, fmt.Errorf("GetOfflineClients json unmarshal error: %w", err)
	}
	var macs []string
	for _, client := range clients {
		if client.Name == "" {
			macs = append(macs, client.Mac)
		}
	}
	return macs, nil
}

func (p *Prune) RemoveOfflineClients() error {
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Prune struct {
	model    string
	ip       string
	port     string
	user     string
	password string

	httpRequest       *HttpRequest
	loginPath         string
	logoutPath        string
	clientHistoryPath string
	cmdRemovalPath    string
}

type Client struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
}

type GetOfflineClientsRsp []*Client

type RemoveOfflineClientsRsp struct {
	Meta map[string]string `json:"meta"`
}

func NewPrune(model, ip, port, user, password string) *Prune {
	var loginPath string
	var logoutPath string
	var clientHistoryPath string
	var cmdRemovalPath string
	if model == ModelController {
		loginPath = "api/login"
		logoutPath = "api/logout"
		clientHistoryPath = "v2/api/site/default/clients/history"
		cmdRemovalPath = "api/s/default/cmd/stamgr"
	} else {
		loginPath = "api/auth/login"
		logoutPath = "api/auth/logout"
		clientHistoryPath = "proxy/network/v2/api/site/default/clients/history"
		cmdRemovalPath = "proxy/network/api/s/default/cmd/stamgr"
	}
	if port == "" {
		port = "443"
	}
	prune := &Prune{
		model:    model,
		ip:       ip,
		port:     port,
		user:     user,
		password: password,

		httpRequest:       NewHttpRequest(30),
		loginPath:         loginPath,
		logoutPath:        logoutPath,
		clientHistoryPath: clientHistoryPath,
		cmdRemovalPath:    cmdRemovalPath,
	}
	return prune
}

func (p *Prune) Run() error {
	fmt.Println("logging in")
	err := p.Login()
	if err != nil {
		return err
	}
	fmt.Println("logged in successfully")
	macs, err := p.GetOfflineClients()
	if err != nil {
		return err
	}
	fmt.Printf("totally %d offline clients\n", len(macs))
	err = p.RemoveOfflineClients(macs)
	if err != nil {
		return err
	}
	err = p.Logout()
	if err != nil {
		return err
	}
	fmt.Println("logged out successfully")
	return nil
}

func (p *Prune) Login() error {
	params := map[string]interface{}{
		"username": p.user,
		"password": p.password,
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	url := fmt.Sprintf("https://%s:%s/%s", p.ip, p.port, p.loginPath)

	rsp, err := p.httpRequest.RequestRaw(url, "POST", params, headers)
	if err != nil {
		return err
	}
	p.httpRequest.StoreCookies(rsp.Cookies())
	p.httpRequest.StoreHeaders(rsp.Header)
	return nil
}

func (p *Prune) Logout() error {
	headers := make(map[string]string)
	token, ok := p.httpRequest.headers["X-Csrf-Token"]
	if ok && len(token) > 0 {
		headers["X-Csrf-Token"] = token[0]
	}
	url := fmt.Sprintf("https://%s:%s/%s", p.ip, p.port, p.logoutPath)

	_, err := p.httpRequest.RequestRaw(url, "POST", nil, headers)
	if err != nil {
		return err
	}
	return nil
}

func (p *Prune) GetOfflineClients() ([]string, error) {
	params := "onlyNonBlocked=true&includeUnifiDevices=true&withinHours=0"
	url := fmt.Sprintf("https://%s:%s/%s?%s", p.ip, p.port, p.clientHistoryPath, params)
	rspStr, err := p.httpRequest.Request(url, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	var rsp GetOfflineClientsRsp
	err = json.Unmarshal(rspStr, &rsp)
	if err != nil {
		return nil, fmt.Errorf("GetOfflineClients json unmarshal error: %w", err)
	}
	var macs []string
	for _, client := range rsp {
		if client.Name == "" {
			macs = append(macs, client.Mac)
		}
	}
	return macs, nil
}

func (p *Prune) RemoveOfflineClients(macs []string) error {
	limit := 5
	start := 0
	end := start + limit
	length := len(macs)
	for {
		if start >= length {
			break
		}
		if end > length {
			end = length
		}
		fmt.Printf("removing clients: %s\n", strings.Join(macs[start:end], " "))
		params := map[string]interface{}{
			"macs": macs,
			"cmd":  "forget-sta",
		}
		headers := make(map[string]string)
		token, ok := p.httpRequest.headers["X-Csrf-Token"]
		if ok && len(token) > 0 {
			headers["X-Csrf-Token"] = token[0]
		}

		url := fmt.Sprintf("https://%s:%s/%s", p.ip, p.port, p.cmdRemovalPath)
		rspStr, err := p.httpRequest.Request(url, "POST", params, headers)
		if err != nil {
			panic(fmt.Errorf("failed to remove offline clients: %w", err))
		}
		var rsp RemoveOfflineClientsRsp
		err = json.Unmarshal(rspStr, &rsp)
		if err != nil {
			panic(fmt.Errorf("RemoveOfflineClients json unmarshal error: %w", err))
		}
		result, ok := rsp.Meta["rc"]
		if ok && result == "ok" {
			fmt.Println("removed successfully")
		}

		start += limit
		end += limit
	}
	return nil
}

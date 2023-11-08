package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"unifi-prune/config"
)

type Unifi struct {
	model    string
	ip       string
	port     string
	user     string
	password string

	HttpRequest       *HttpRequest
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

func NewUnifi(model, ip, port, user, password string) *Unifi {
	var loginPath string
	var logoutPath string
	var clientHistoryPath string
	var cmdRemovalPath string
	if model == config.ModelController {
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
	prune := &Unifi{
		model:    model,
		ip:       ip,
		port:     port,
		user:     user,
		password: password,

		HttpRequest:       NewHttpRequest(),
		loginPath:         loginPath,
		logoutPath:        logoutPath,
		clientHistoryPath: clientHistoryPath,
		cmdRemovalPath:    cmdRemovalPath,
	}
	return prune
}

func (u *Unifi) Login() error {
	params := map[string]interface{}{
		"username": u.user,
		"password": u.password,
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	url := fmt.Sprintf("https://%s:%s/%s", u.ip, u.port, u.loginPath)

	rsp, err := u.HttpRequest.RequestRaw(url, "POST", params, headers)
	if err != nil {
		return err
	}
	u.HttpRequest.StoreCookies(rsp.Cookies())
	u.HttpRequest.StoreHeaders(rsp.Header)
	return nil
}

func (u *Unifi) Logout() error {
	headers := make(map[string]string)
	token, ok := u.HttpRequest.Headers["X-Csrf-Token"]
	if ok && len(token) > 0 {
		headers["X-Csrf-Token"] = token[0]
	}
	url := fmt.Sprintf("https://%s:%s/%s", u.ip, u.port, u.logoutPath)

	_, err := u.HttpRequest.RequestRaw(url, "POST", nil, headers)
	if err != nil {
		return err
	}
	return nil
}

func (u *Unifi) GetOfflineClients() ([]string, error) {
	params := "onlyNonBlocked=true&includeUnifiDevices=true&withinHours=0"
	url := fmt.Sprintf("https://%s:%s/%s?%s", u.ip, u.port, u.clientHistoryPath, params)
	rspStr, err := u.HttpRequest.Request(url, "GET", nil, nil)
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

func (u *Unifi) RemoveOfflineClients(macs []string) error {
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
		token, ok := u.HttpRequest.Headers["X-Csrf-Token"]
		if ok && len(token) > 0 {
			headers["X-Csrf-Token"] = token[0]
		}

		url := fmt.Sprintf("https://%s:%s/%s", u.ip, u.port, u.cmdRemovalPath)
		rspStr, err := u.HttpRequest.Request(url, "POST", params, headers)
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

func (u *Unifi) Prune() error {
	fmt.Println("logging in")
	err := u.Login()
	if err != nil {
		return err
	}
	fmt.Println("logged in successfully")
	macs, err := u.GetOfflineClients()
	if err != nil {
		return err
	}
	fmt.Printf("totally %d offline clients\n", len(macs))
	err = u.RemoveOfflineClients(macs)
	if err != nil {
		return err
	}
	err = u.Logout()
	if err != nil {
		return err
	}
	fmt.Println("logged out successfully")
	return nil
}

func (u *Unifi) Recover() {
	if r := recover(); r != nil {
		fmt.Println(r)
		if len(u.HttpRequest.Cookies) == 0 {
			return
		}
		err := u.Logout()
		if err != nil {
			fmt.Printf("failed to logout: %s\n", err.Error())
		}
		fmt.Println("logged out successfully")
	}
}

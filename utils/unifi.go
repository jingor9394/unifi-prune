package utils

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unifi-prune/configs"
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
	clientActivePath  string
	clientHistoryPath string
	cmdRemovalPath    string
	wlanConfigPath    string
}

type Client struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
}

type WlanConfig struct {
	Config struct {
		Name          string   `json:"name"`
		MacFilterList []string `json:"mac_filter_list"`
	} `json:"configuration"`
}

type GetClientsRsp []*Client
type GetWlanConfig []*WlanConfig

type RemoveOfflineClientsRsp struct {
	Meta map[string]string `json:"meta"`
}

func NewUnifi(model, ip, port, user, password string) *Unifi {
	var loginPath string
	var logoutPath string
	var clientActivePath string
	var clientHistoryPath string
	var cmdRemovalPath string
	var wlanConfigPath string
	if model == configs.ModelController {
		loginPath = configs.ControllerLoginPath
		logoutPath = configs.ControllerLogoutPath
		clientActivePath = configs.ControllerClientActivePath
		clientHistoryPath = configs.ControllerClientHistoryPath
		cmdRemovalPath = configs.ControllerCmdRemovalPath
		wlanConfigPath = configs.ControllerWlanConfigPath
	} else {
		loginPath = configs.LoginPath
		logoutPath = configs.LogoutPath
		clientActivePath = configs.ClientActivePath
		clientHistoryPath = configs.ClientHistoryPath
		cmdRemovalPath = configs.CmdRemovalPath
		wlanConfigPath = configs.WlanConfigPath
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
		clientActivePath:  clientActivePath,
		clientHistoryPath: clientHistoryPath,
		cmdRemovalPath:    cmdRemovalPath,
		wlanConfigPath:    wlanConfigPath,
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

func (u *Unifi) GetOfflineClients() ([]*Client, error) {
	params := "onlyNonBlocked=true&includeUnifiDevices=true&withinHours=0"
	url := fmt.Sprintf("https://%s:%s/%s?%s", u.ip, u.port, u.clientHistoryPath, params)
	rspStr, err := u.HttpRequest.Request(url, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	var rsp GetClientsRsp
	err = json.Unmarshal(rspStr, &rsp)
	if err != nil {
		return nil, fmt.Errorf("GetOfflineClients json unmarshal error: %w", err)
	}
	return rsp, nil
}

func (u *Unifi) GetActiveClients() ([]*Client, error) {
	params := "includeTrafficUsage=true&includeUnifiDevices=true"
	url := fmt.Sprintf("https://%s:%s/%s?%s", u.ip, u.port, u.clientActivePath, params)
	rspStr, err := u.HttpRequest.Request(url, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	var rsp GetClientsRsp
	err = json.Unmarshal(rspStr, &rsp)
	if err != nil {
		return nil, fmt.Errorf("GetActiveClients json unmarshal error: %w", err)
	}
	return rsp, nil
}

func (u *Unifi) GetRemovedMacs(clients []*Client) []string {
	var macs []string
	for _, client := range clients {
		if client.Name == "" {
			macs = append(macs, client.Mac)
		}
	}
	fmt.Printf("totally %d offline clients: %s\n", len(macs), strings.Join(macs, ", "))
	return macs
}

func (u *Unifi) GetClientsMap() (map[string]string, error) {
	clientsMap := make(map[string]string)
	activeClients, err := u.GetActiveClients()
	if err != nil {
		return nil, err
	}
	offlineClients, err := u.GetOfflineClients()
	if err != nil {
		return nil, err
	}
	clients := append(activeClients, offlineClients...)
	for _, client := range clients {
		clientsMap[client.Mac] = client.Name
	}
	return clientsMap, nil
}

func (u *Unifi) GetWlanConfigs() ([]*WlanConfig, error) {
	url := fmt.Sprintf("https://%s:%s/%s", u.ip, u.port, u.wlanConfigPath)
	rspStr, err := u.HttpRequest.Request(url, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	var rsp GetWlanConfig
	err = json.Unmarshal(rspStr, &rsp)
	if err != nil {
		return nil, fmt.Errorf("GetWlanConfigs json unmarshal error: %w", err)
	}
	return rsp, nil
}

func (u *Unifi) RemoveOfflineMacs(macs []string) error {
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
			panic(fmt.Errorf("RemoveOfflineMacs json unmarshal error: %w", err))
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

func (u *Unifi) Prune(dryRun bool) error {
	err := u.Login()
	if err != nil {
		return err
	}
	fmt.Println("logged in successfully")
	clients, err := u.GetOfflineClients()
	if err != nil {
		return err
	}
	macs := u.GetRemovedMacs(clients)
	if !dryRun {
		err = u.RemoveOfflineMacs(macs)
		if err != nil {
			return err
		}
	}
	err = u.Logout()
	if err != nil {
		return err
	}
	fmt.Println("logged out successfully")
	return nil
}

func (u *Unifi) PrintMacFilterList() error {
	err := u.Login()
	if err != nil {
		return err
	}
	fmt.Println("logged in successfully")

	clientsMap, err := u.GetClientsMap()
	if err != nil {
		return err
	}
	wlanConfigs, err := u.GetWlanConfigs()
	if err != nil {
		return err
	}
	for _, wlanConfig := range wlanConfigs {
		fmt.Printf("[%s]\n", wlanConfig.Config.Name)
		// Sort keys first by value and then iterate map
		macFilterMap := make(map[string]string)
		macs := make([]string, 0, len(wlanConfig.Config.MacFilterList))
		for _, mac := range wlanConfig.Config.MacFilterList {
			macs = append(macs, mac)
			if name, ok := clientsMap[mac]; ok {
				macFilterMap[mac] = name
				continue
			}
			macFilterMap[mac] = "Unknown"
		}
		sort.Slice(macs, func(i, j int) bool {
			return macFilterMap[macs[i]] < macFilterMap[macs[j]]
		})
		for _, mac := range macs {
			fmt.Printf("%s: %s\n", mac, macFilterMap[mac])
		}
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

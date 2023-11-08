package utils

import (
	"errors"
	"flag"
	"fmt"
	"unifi-prune/config"
)

type Params struct {
	Model    string
	IP       string
	Port     string
	User     string
	Password string
}

var InvalidParams = errors.New("invalid params")

func NewParams() (*Params, error) {
	var model string
	var ip string
	var user string
	var password string
	var port string

	flag.StringVar(&model, "m", "", "")
	flag.StringVar(&ip, "g", "", "")
	flag.StringVar(&user, "u", "", "")
	flag.StringVar(&port, "p", "", "")
	flag.Usage = func() {
		usage := `Usages:
-m <UDMPro/UDR/Controller>
    Unifi model
-g
    Unifi console ip address
-u
    Unifi console user
-p
    Unifi console password`
		fmt.Println(usage)
	}
	flag.Parse()

	if model == "" || ip == "" || user == "" {
		flag.Usage()
		return nil, InvalidParams
	}
	models := map[string]bool{
		config.ModelUDMPro:     true,
		config.ModelUDR:        true,
		config.ModelController: true,
	}
	if _, ok := models[model]; !ok {
		flag.Usage()
		return nil, InvalidParams
	}

	fmt.Println("please enter password:")
	_, err := fmt.Scanln(&password)
	if err != nil {
		fmt.Println(err)
		return nil, InvalidParams
	}

	params := &Params{
		Model:    model,
		IP:       ip,
		Port:     port,
		User:     user,
		Password: password,
	}
	return params, nil
}

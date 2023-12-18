package utils

import (
	"errors"
	"flag"
	"fmt"
	"unifi-prune/configs"
)

type Params struct {
	Model    string
	IP       string
	Port     string
	User     string
	Password string
	DryRun   bool
}

var InvalidParams = errors.New("invalid params")

func NewParams() (*Params, error) {
	var model string
	var ip string
	var user string
	var password string
	var port string
	var dryRun bool

	flag.StringVar(&model, "m", "", "")
	flag.StringVar(&ip, "g", "", "")
	flag.StringVar(&user, "u", "", "")
	flag.StringVar(&port, "p", "", "")
	flag.BoolVar(&dryRun, "d", false, "")
	flag.Usage = func() {
		usage := `Usages:
-m <Console/Controller>
    Unifi model
-g
    Unifi console ip address
-p
    Unifi console port (Optional)
-u
    Unifi console user
-d <true>
    Dry run`
		fmt.Println(usage)
	}
	flag.Parse()

	if model == "" || ip == "" || user == "" {
		flag.Usage()
		return nil, InvalidParams
	}
	models := map[string]bool{
		configs.ModelConsole:    true,
		configs.ModelController: true,
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
		DryRun:   dryRun,
	}
	return params, nil
}

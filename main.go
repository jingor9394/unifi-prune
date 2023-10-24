package main

import (
	"flag"
	"fmt"
)

const (
	ModelUDMPro     = "UDMPro"
	ModelUDR        = "UDR"
	ModelController = "Controller"
)

func main() {
	var model string
	var ip string
	var user string
	var password string

	flag.StringVar(&model, "m", "", "")
	flag.StringVar(&ip, "g", "", "")
	flag.StringVar(&user, "u", "", "")
	flag.StringVar(&password, "p", "", "")
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

	if model == "" || ip == "" || user == "" || password == "" {
		flag.Usage()
		return
	}
	models := map[string]bool{
		ModelUDMPro:     true,
		ModelUDR:        true,
		ModelController: true,
	}
	if _, ok := models[model]; !ok {
		flag.Usage()
		return
	}

	prune := NewPrune(model, ip, user, password)
	err := prune.Login()
	if err != nil {
		fmt.Println(err)
	}
}

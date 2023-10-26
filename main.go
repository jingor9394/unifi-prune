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
	var port string

	flag.StringVar(&model, "m", "", "")
	flag.StringVar(&ip, "g", "", "")
	flag.StringVar(&user, "u", "", "")
	flag.StringVar(&password, "p", "", "")
	flag.StringVar(&port, "P", "", "")
	flag.Usage = func() {
		usage := `Usages:
-m <UDMPro/UDR/Controller>
    Unifi model
-g
    Unifi console ip address
-u
    Unifi console user
-p
    Unifi console password
-P
    Unifi console port, default is 443`
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

	run(model, ip, port, user, password)
}

func run(model, ip, port, user, password string) {
	prune := NewPrune(model, ip, port, user, password)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			if len(prune.httpRequest.cookies) == 0 {
				return
			}
			err := prune.Logout()
			if err != nil {
				fmt.Printf("failed to logout: %s\n", err.Error())
			}
			fmt.Println("logged out successfully")
		}
	}()
	err := prune.Run()
	if err != nil {
		fmt.Println(err)
	}
}

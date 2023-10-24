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

type UnifiPrune struct {
	Model    string
	Ip       string
	User     string
	Password string
	URL      string
}

func NewUnifiPrune(model, ip, user, password string) *UnifiPrune {
	unifiPrune := &UnifiPrune{
		Model:    model,
		Ip:       ip,
		User:     user,
		Password: password,
	}
	if model == ModelUDMPro {

	} else {

	}
	return unifiPrune
}

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

	fmt.Println(model)
	fmt.Println(ip)
	fmt.Println(user)
	fmt.Println(password)
}

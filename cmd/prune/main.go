package main

import (
	"flag"
	"fmt"
	"unifi-prune/config"
	"unifi-prune/utils"
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
		return
	}
	models := map[string]bool{
		config.ModelUDMPro:     true,
		config.ModelUDR:        true,
		config.ModelController: true,
	}
	if _, ok := models[model]; !ok {
		flag.Usage()
		return
	}

	fmt.Println("please enter password:")
	_, err := fmt.Scanln(&password)
	if err != nil {
		fmt.Println(err)
	}

	run(model, ip, port, user, password)
}

func run(model, ip, port, user, password string) {
	unifi := utils.NewUnifi(model, ip, port, user, password)
	defer unifi.Recover()
	err := unifi.Prune()
	if err != nil {
		fmt.Println(err)
	}
}

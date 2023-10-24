package main

import (
	"flag"
	"fmt"
)

func main() {
	var model string
	var gateway string
	var user string
	var password string

	flag.StringVar(&model, "m", "", "")
	flag.StringVar(&gateway, "g", "", "")
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

	if model == "" || gateway == "" || user == "" || password == "" {
		flag.Usage()
		return
	}
	models := map[string]bool{
		"UDMPro":     true,
		"UDR":        true,
		"Controller": true,
	}
	if _, ok := models[model]; !ok {
		flag.Usage()
		return
	}

	fmt.Println(model)
	fmt.Println(gateway)
	fmt.Println(user)
	fmt.Println(password)
}

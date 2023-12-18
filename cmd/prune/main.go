package main

import (
	"fmt"
	"unifi-prune/utils"
)

func main() {
	params, err := utils.NewParams()
	if err != nil {
		return
	}
	unifi := utils.NewUnifi(params.Model, params.IP, params.Port, params.User, params.Password)
	defer unifi.Recover()
	err = unifi.Prune(params.DryRun)
	if err != nil {
		fmt.Println(err)
	}
}

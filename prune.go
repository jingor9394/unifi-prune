package main

import "fmt"

type Prune struct {
	model    string
	ip       string
	user     string
	password string
	url      string

	httpRequest *HttpRequest
}

func NewPrune(model, ip, user, password string) *Prune {
	prune := &Prune{
		model:       model,
		ip:          ip,
		user:        user,
		password:    password,
		httpRequest: NewHttpRequest(10),
	}
	if model == ModelUDMPro {

	} else {

	}
	return prune
}

func (p *Prune) Login() error {
	data := map[string]interface{}{
		"username": p.user,
		"password": p.password,
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	url := fmt.Sprintf("http://%s/api/auth/login", p.ip)
	rsp, err := p.httpRequest.Request(url, "POST", data, headers)
	if err != nil {
		return err
	}
	fmt.Println(string(rsp))
	return nil
}

package main

import (
	"fmt"
)

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

func (p *Prune) Run() error {
	err := p.Login()
	if err != nil {
		return err
	}
	err = p.Logout()
	if err != nil {
		return err
	}
	return nil
}

func (p *Prune) Login() error {
	data := map[string]interface{}{
		"username": p.user,
		"password": p.password,
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	url := fmt.Sprintf("https://%s/api/auth/login", p.ip)

	_, err := p.httpRequest.Request(url, "POST", data, headers)
	if err != nil {
		return err
	}
	p.httpRequest.GetCookies()
	return nil
}

func (p *Prune) Logout() error {
	url := fmt.Sprintf("https://%s/api/auth/logout", p.ip)
	headers := make(map[string]string)
	token, ok := p.httpRequest.headers["X-Csrf-Token"]
	if ok && len(token) > 0 {
		headers["X-Csrf-Token"] = token[0]
	}

	_, err := p.httpRequest.Request(url, "POST", nil, headers)
	if err != nil {
		return err
	}
	return nil
}

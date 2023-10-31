package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type HttpRequest struct {
	timeout int
	cookies []*http.Cookie
	headers http.Header
}

func NewHttpRequest(timeout int) *HttpRequest {
	if timeout == 0 {
		timeout = 10
	}
	httpRequest := &HttpRequest{
		timeout: timeout,
	}
	return httpRequest
}

func (r *HttpRequest) getRequest(url, method string, params map[string]interface{}) (*http.Request, error) {
	var paramStr string
	if params != nil {
		jsonStr, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("getRequest json marshal error: %w", err)
		}
		paramStr = string(jsonStr)
	}
	reqParams := strings.NewReader(paramStr)
	req, err := http.NewRequest(method, url, reqParams)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *HttpRequest) setHeader(req *http.Request, headers map[string]string) {
	if headers != nil {
		for header, val := range headers {
			if header == "Host" {
				req.Host = val
			} else {
				req.Header.Add(header, val)
			}
		}
	}
}

func (r *HttpRequest) getClient(req *http.Request) *http.Client {
	reqTimeout := time.Duration(r.timeout) * time.Second
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		},
	}
	jar, _ := cookiejar.New(nil)
	if len(r.cookies) != 0 {
		jar.SetCookies(req.URL, r.cookies)
	}
	client := &http.Client{
		Timeout:   reqTimeout,
		Transport: tr,
		Jar:       jar,
	}
	return client
}

func (r *HttpRequest) getResponseString(rsp *http.Response) ([]byte, error) {
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code: %d", rsp.StatusCode)
	}
	rspStr, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return rspStr, nil
}

func (r *HttpRequest) RequestRaw(url, method string, params map[string]interface{}, headers map[string]string) (*http.Response, error) {
	req, err := r.getRequest(url, method, params)
	if err != nil {
		return nil, err
	}
	r.setHeader(req, headers)

	client := r.getClient(req)
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	return rsp, nil

}

func (r *HttpRequest) Request(url, method string, params map[string]interface{}, headers map[string]string) ([]byte, error) {
	rsp, err := r.RequestRaw(url, method, params, headers)
	if err != nil {
		return nil, err
	}

	rspStr, err := r.getResponseString(rsp)
	if err != nil {
		return nil, err
	}
	return rspStr, nil
}

func (r *HttpRequest) SetCookies(cookies []*http.Cookie) {
	r.cookies = cookies
}

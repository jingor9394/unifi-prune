package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpRequest struct {
	timeout int
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

func (r *HttpRequest) Request(url, method string, data map[string]interface{}, headers map[string]string) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	reqData := strings.NewReader(string(dataStr))
	req, err := http.NewRequest(method, url, reqData)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for header, val := range headers {
			if header == "Host" {
				req.Host = val
			} else {
				req.Header.Add(header, val)
			}
		}
	}
	reqTimeout := time.Duration(r.timeout) * time.Second
	client := &http.Client{Timeout: reqTimeout}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	rspStr, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return rspStr, nil
}

package main

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"crypto/tls"
)

// the goal here is to build
// 	- tcp "tcp is the network that send the data" server
// 	- for http "is the protocol for sending these data"

const (
	HostDevelopment = "https://api.development.push.apple.com"
)

type Client struct {
	HttpClient  *http.Client
	Certificate tls.Certificate
	Host        string
}

func NewClient(certificate tls.Certificate) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &Client{
		HttpClient:  &http.Client{Transport: transport},
		Certificate: certificate,
	}
}

func setHeaders(r *http.Request, n *Notification) {
	r.Header.Set("Content-Type", "application/json")
	if n.Topic != "" {
		r.Header.Set("apns-topic", n.Topic)
	}
	if n.ApnsID != "" {
		r.Header.Set("apns-id", n.ApnsID)
	}

}

func (c *Client) Push(n *Notification , host string) (*Response, error) {
	url := fmt.Sprintf("%v/3/device/%v", host, n.DeviceToken)
	// convert the payload to ARRAY OF BYTES to be handled correctly
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(n.Payload.([]byte)))
	setHeaders(req, n)
	httpRes, httpErr := c.HttpClient.Do(req)

	if httpErr != nil {
		return nil, httpErr
	}
	defer httpRes.Body.Close()

	res := &Response{}
	res.StatusCode = httpRes.StatusCode
	res.ApnsID = httpRes.Header.Get("apns-id")
	if res.StatusCode == http.StatusOK {
		return res, nil
	} else {
		err := &APNSError{}
		json.NewDecoder(httpRes.Body).Decode(err)
		return res, err
	}
}

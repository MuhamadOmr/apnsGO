package main

import (
	"fmt"
	"net/http"
	"bytes"
	"golang.org/x/net/http2"
	"crypto/tls"
	"net"
	"time"
	"sync"
)


const (
	HostDevelopment = "https://api.development.push.apple.com"
)

type HttpsClient struct {
	Client *http.Client

}



var DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
	dialer := &net.Dialer{
		//Connection timeout with certificate
		Timeout: 20 * time.Second,
		//Duration of activity after not using the connection
		KeepAlive: 60 * time.Second,
	}
	return tls.DialWithDialer(dialer, network, addr, cfg)
}


func NewClient(certificate tls.Certificate) *HttpsClient {

	newHttpsClient := &HttpsClient{}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	tlsConfig.BuildNameToCertificate()

	newHttpsClient.Client = &http.Client{
		//HTTP2
		Transport: &http2.Transport{
			TLSClientConfig: tlsConfig,
			DialTLS:         DialTLS,
		},
	}

	return newHttpsClient
}

func setHeaders(r *http.Request, n *Notification) {
	r.Header.Set("Content-Type", "application/json")
	if n.ApnsID != "" {
		r.Header.Set("apns-id", n.ApnsID)
	}

}

func (c *HttpsClient) Push(n *Notification ,wg *sync.WaitGroup , host string , DeviceToken string) {
	defer wg.Done()

	url := fmt.Sprintf("%v/3/device/%v", host, DeviceToken)
	// convert the payload to ARRAY OF BYTES to be handled correctly
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(n.Payload.([]byte)))
	setHeaders(req, n)
	httpRes, httpErr := c.Client.Do(req)

	if httpErr != nil {
		fmt.Println(httpErr)
	}

	if httpRes != nil {
		defer httpRes.Body.Close()

		fmt.Println(httpRes.Status)
	}

}


func Act( Client *HttpsClient , n *Notification, host string , DeviceTokens []string) {

	var wg sync.WaitGroup

	for i := 0; i < len(DeviceTokens); i++ {
		wg.Add(1)
		go Client.Push(n , &wg, host ,DeviceTokens[i])
	}
	wg.Wait()

}
package main


import (
	"fmt"
	"log"
)

const cert = "CERT HERE"
const key = "KEY HERE"

func main() {

	byteCert := []byte(cert)
	// generate the certificate
	tlsCert,err := FromPemBytes(byteCert , key)

	if err != nil {
		log.Println("Certificate Error: ", err)
		return
	}

	//make the notification
	notification := &Notification{}

	notification.DeviceToken = "TOKEN HERE"
	notification.Topic = "com.sideshow.Apns2"
	notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)

	// push the notification and handle the response


	client := NewClient(tlsCert)
	res, err := client.Push(notification , HostDevelopment )

	if err != nil {
		log.Println("APNS Error: ", err)
		return
	}



	fmt.Println(res);

}
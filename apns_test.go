package apns

import (
	"apns"
	"encoding/base64"
	"encoding/hex"
	"log"
	// "time"
	// "sync"
)

/*
openssl pkcs12 -in aps-dev.p12 -out crt.pem -clcerts -nokeys
openssl pkcs12 -in aps-dev.p12 -out key.pem -nocerts -nodes
*/

func Test() {
	certPath := "/Users/xuhua/Documents/keys/card"
	s := "M1ieMadAhKtBSBoGWxhcHR8M2qOJdMg0iWq0vwgjQGM="
	token, err := base64.StdEncoding.DecodeString(s)
	log.Println(err, len(token))

	client := apns.NewClient(apns.ApnsSandboxGateway, certPath+"/crt.pem", certPath+"/key.pem")

	dict := apns.NewAlertDictionary()
	dict.LocKey = "TEST"
	dict.LocArgs = []string{"11"}
	payload := apns.NewPayload()
	payload.Alert = dict
	payload.Badge = 18
	// payload.Sound = "bingbong.aiff"
	pn := apns.NewNotification()
	pn.DeviceToken = hex.EncodeToString(token)
	// pn.Identifier = 10
	pn.AddPayload(payload)
	client.Append(pn)
	client.Send()

}

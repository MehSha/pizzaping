package main

import (
	"fmt"

	webpush "github.com/SherClockHolmes/webpush-go"
)

const (
	testSubs        = `{"endpoint":"https://fcm.googleapis.com/fcm/send/cHqYNCAwmG8:APA91bFATeUwX3v7ZsVDwiOabOn_jGW--2Bi3nxzyRsFr3dG5WDaMJZEr11ie4MmS489qQnd3XCl0MdeKPov5AN6aEFjiC3roCsd2n86K0W8pz2Ixak_KRj_tpwp-nVOfO6cOlYW8xrF","expirationTime":null,"keys":{"p256dh":"BBcu7R_9GIzXvMZib0oDQjTJn9TpLvb2vM1Ai6UH4mESeueHYo1jleyijt0AG36UNg-r9Heisz2twtJ5z7AJRww","auth":"GOTTP1ySyPV-meJNaL9mKQ"}}`
	vapidPublicKey  = "BElNerWtL4YzaBkQdi_FwLrKbzqmXDDipN4FhyEzXLn9-Gm9eLh5Midcsfi51xqlECG1cmXS-1kyx7ZN4NK3lqY"
	vapidPrivateKey = "6gtifsr-nw3fs_kSrAuStdUxA2g029GA6LSMqs6v5hc"
)

func pushNotif(rest Restaurant, message string) {
	// Decode subscription
	s := &webpush.Subscription{
		Endpoint: rest.EndPoint,
		Keys: webpush.Keys{
			P256dh: rest.P256DH,
			Auth:   rest.Auth,
		},
	}

	// Send Notification
	_, err := webpush.SendNotification([]byte(message), s, &webpush.Options{
		Subscriber:      "example@example.com", // Do not include "mailto:"
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})
	if err != nil {
		fmt.Println("error sending notification:", err)
	}
}

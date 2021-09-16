package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "wss", Host: "127.0.0.1:8020", Path: "/piledriver"}

	dialer := &websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	dial, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer dial.Close()

	for {
		var cmd string
		fmt.Print("exe: ")
		fmt.Scanln(&cmd)
		if cmd == "" || cmd == "exit" {
			break
		}

		err := dial.WriteMessage(websocket.TextMessage, []byte(cmd))
		if err != nil {
			log.Println(err)
			return
		}

		_, p, err := dial.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println("r: ", string(p))
	}
}

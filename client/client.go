package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var localHost = flag.String("local_host", "127.0.0.1:9871", "Local Host")
var localUser = flag.String("local_user", "piledriver", "Local User")
var localPassword = flag.String("local_password", "piledriver", "Local Password")
var pacDns = flag.String("pac_dns", "", "Pac Dns")
var remoteHost = flag.String("remote_host", "127.0.0.1:8020", "Remote Host")
var remotePath = flag.String("remote_path", "/piledriver", "Remote PATH")
var token = flag.String("token", "", "token auth")

func main() {
	flag.Parse()

	// Local
	addr, err := net.ResolveTCPAddr("tcp", *localHost)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		accept, err := conn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		c, err := initClient()
		if err != nil {
			log.Println(err)
			continue
		}
		go c.accept(accept)
	}
}

type client struct {
	conn *websocket.Conn
}

func (c *client) accept(conn net.Conn) {

}

func initClient() (*client, error) {
	u := url.URL{Scheme: "wss", Host: *remoteHost, Path: *remotePath}

	dialer := &websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	dial, _, err := dialer.Dial(u.String(), http.Header{"token": []string{*token}})
	if err != nil {
		return nil, err
	}

	return &client{conn: dial}, nil
	//for {
	//	var cmd string
	//	fmt.Print("exe: ")
	//	fmt.Scanln(&cmd)
	//	if cmd == "" || cmd == "exit" {
	//		break
	//	}
	//
	//	err := dial.WriteMessage(websocket.TextMessage, []byte(cmd))
	//	if err != nil {
	//		log.Println(err)
	//		return nil, err
	//	}
	//
	//	_, p, err := dial.ReadMessage()
	//	if err != nil {
	//		log.Println(err)
	//		break
	//	}
	//
	//	fmt.Println("r: ", string(p))
	//}
	//
	//return nil, err
}

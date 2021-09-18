package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

var localHost = flag.String("local_host", "127.0.0.1:9871", "Local Host")
var localUser = flag.String("local_user", "piledriver", "Local User")
var localPassword = flag.String("local_password", "piledriver", "Local Password")
var pacDns = flag.String("pac_dns", "", "Pac Dns")
var remoteHost = flag.String("remote_host", "127.0.0.1:8020", "Remote Host")
var remotePath = flag.String("remote_path", "/piledriver", "Remote PATH")
var token = flag.String("token", "piledriver", "token auth")

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Llongfile)

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
			//log.Println(err)
			continue
		}
		go c.accept(accept)
	}
}

type client struct {
	conn *websocket.Conn
}

func (c *client) accept(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	if conn == nil {
		return
	}

	defer conn.Close()

	var b [1024]byte
	_, err := conn.Read(b[:])
	if err != nil {
		return
	}

	if b[0] != 0x05 {
		return
	}

	conn.Write([]byte{0x05, 0x00})
	n, err := conn.Read(b[:])
	if err != nil {
		return
	}

	// 解析目的地
	var host, port string
	switch b[3] {
	case 0x01: // ip
		host = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case 0x03: // domain
		host = string(b[5 : n-2]) //b[4]表示域名的长度
	case 0x04: // ipv6
		return
	}

	port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
	addr := net.JoinHostPort(host, port)

	if _, err := conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); err != nil { //响应客户端连接成功
		return
	}

	err = c.conn.WriteMessage(websocket.TextMessage, []byte(addr))
	if err != nil {
		log.Println(err)
		return
	}

	go copy1(c.conn, conn)
	copy2(conn, c.conn)
}

func initClient() (*client, error) {
	u := url.URL{Scheme: "wss", Host: *remoteHost, Path: *remotePath}

	dialer := &websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	dial, _, err := dialer.Dial(u.String(), http.Header{"token": []string{*token}})
	if err != nil {
		return nil, err
	}

	return &client{conn: dial}, nil
}

func copy1(conn *websocket.Conn, client io.Reader) {
	for {
		var b [1024]byte
		read, err := client.Read(b[:])
		if err != nil {
			if err == io.EOF {
				conn.Close()
				break
			}
			break
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, b[:read]); err != nil {
			break
		}
	}
}

func copy2(client io.Writer, conn *websocket.Conn) {
	for {
		_, byt, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if _, err := client.Write(byt); err != nil {
			break
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/dollarkillerx/urllib"
	"github.com/gorilla/websocket"
)

var localHost = flag.String("local_host", "0.0.0.0:443", "Local Host")
var token = flag.String("token", "piledriver", "token auth")
var url = flag.String("url", "piledriver", "url")
var target = flag.String("target", "https://www.dollarkiller.com/light/view/", "target")

type PiledriverHandler struct{}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clientMap sync.Map
)

type Handler func(write http.ResponseWriter, req *http.Request)

func (p *PiledriverHandler) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	switch req.URL.String() {
	case fmt.Sprintf("/%s", *url):
		if req.Header.Get("token") != *token {
			write.WriteHeader(401)
			write.Write([]byte("401"))
			log.Println("RemoteAddr: ", req.RemoteAddr, " token : ", *token, "  error")
			return
		}

		conn, err := upgrader.Upgrade(write, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			if err != websocket.ErrCloseSent {
				//log.Println(err)
			}
			return
		}

		addr := string(data)
		ipAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Println(err)
			return
		}
		dial, err := net.DialTCP("tcp", nil, ipAddr)
		if err != nil {
			log.Println(err)
			return
		}

		go copy1(dial, conn)
		copy2(conn, dial)
	default:
		targetUr := fmt.Sprintf("%s/%s", *target, req.URL.String())
		code, original, err := urllib.Get(targetUr).ByteOriginal()
		if err != nil {
			write.WriteHeader(500)
			write.Write([]byte(err.Error()))
			return
		}
		write.WriteHeader(code)

		switch {
		case strings.Contains(req.URL.String(), ".css"):
			write.Header().Set("content-type", "text/css; charset=utf-8")
		case strings.Contains(req.URL.String(), ".js"):
			write.Header().Set("content-type", "application/javascript; charset=utf-8")
		default:
			write.Header().Set("content-type", "text/html charset=utf-8")
		}
		write.Write(original)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	ser := &PiledriverHandler{}
	err := http.ListenAndServeTLS(*localHost, "./key/server.crt", "./key/server.key", ser)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func copy1(server io.Writer, conn *websocket.Conn) {
	for {
		_, r, err := conn.ReadMessage()
		if err != nil {
			return
		}

		_, err = server.Write(r)
		if err != nil {
			return
		}
	}
}

func copy2(conn *websocket.Conn, client io.Reader) {
	for {
		var b [1024]byte
		read, err := client.Read(b[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, b[:read]); err != nil {
			break
		}
	}
}

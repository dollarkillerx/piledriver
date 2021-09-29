package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"plumber/pkg/models"
	"strings"
	"sync"

	"github.com/dollarkillerx/urllib"
	"github.com/gorilla/websocket"
)

var localHost = flag.String("local_host", "0.0.0.0:443", "Local Host")
var token = flag.String("token", "piledriver", "token auth")
var url = flag.String("url", "piledriver", "url")
var target = flag.String("target", "https://www.dollarkiller.com/light/view/", "target")
var debug = flag.Bool("debug", false, "debug")

type PiledriverHandler struct {
	clientMap map[string]*net.TCPConn
	mu        sync.Mutex
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
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
			if *debug {
				log.Println(err)
			}
			return
		}
		p.core(conn)
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

func (p *PiledriverHandler) core(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if err != websocket.ErrCloseSent {
				if *debug {
					log.Println(err)
				}
			}
			return
		}

		tml := models.Tml{}
		err = tml.FromBytes(data)
		if err != nil {
			if *debug {
				log.Println(err)
			}
			return
		}

		switch {
		case tml.Start:
			// TODO： 是否通知错误
			addr, err := net.ResolveTCPAddr("tcp", string(tml.Data))
			if err != nil {
				if *debug {
					log.Println(err)
				}
				return
			}
			dial, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				if *debug {
					log.Println(err)
				}
				return
			}
			dial.SetLinger(0)
			p.clientMap[tml.ID] = dial
			continue
		case tml.Close:
			p.clientMap[tml.ID].Close()
			delete(p.clientMap, tml.ID)
			//TODO: 加代码
			continue
		default:
			tcpConn, ex := p.clientMap[tml.ID]
			if ex {
				// TODO: 改
				return
			}

			_, err := tcpConn.Write(tml.Data)
			if err != nil {
				// TODO: 改
				return
			}
		}
	}
}

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	ser := &PiledriverHandler{clientMap: map[string]*net.TCPConn{}}
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

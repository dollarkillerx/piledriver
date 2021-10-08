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

// close 关闭一个细小链接
func (p *PiledriverHandler) close(conn *websocket.Conn, id string) {
	fmt.Println("close ...")
	respTml := models.Tml{
		ID:    id,
		Close: true,
	}

	conn.WriteMessage(websocket.BinaryMessage, respTml.ToBytes())
}

func (p *PiledriverHandler) core(conn *websocket.Conn) {
	fmt.Println("in: ", conn.RemoteAddr())
	defer func() {
		conn.Close()
		fmt.Println("close: ", conn.RemoteAddr())
	}()

	// 解决主要矛盾:
	// 1. 删除不必要的链接
	// 2. 区分通讯链接 和 数据链接

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			if err != websocket.ErrCloseSent {
				if *debug {
					log.Println(err)
				}
			}
			return
		}

		if msgType != websocket.BinaryMessage {
			continue
		}

		tml := models.Tml{}
		err = tml.FromBytes(data)
		if err != nil {
			if *debug {
				log.Println(err)
			}
			return
		}

		if *debug {
			fmt.Println("In Msg: ", tml.ID)
		}

		switch {
		case tml.Start:
			// TODO： 是否通知错误
			addr, err := net.ResolveTCPAddr("tcp", string(tml.Data))
			if err != nil {
				if *debug {
					log.Println(err)
				}
				p.close(conn, tml.ID)
				return
			}
			dial, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				if *debug {
					log.Println(err)
				}
				p.close(conn, tml.ID)
				return
			}
			dial.SetLinger(0)
			p.clientMap[tml.ID] = dial

			go func() {
				// curl --socks5 127.0.0.1:8989 http://www.baidu.com
				defer func() {
					if err := recover(); err != nil {
						fmt.Printf("%s\n", err)
						return
					}
				}()
				for {
					var b [1024]byte
					read, err := dial.Read(b[:])
					if err != nil {
						if *debug {
							log.Println(err)
						}

						p.close(conn, tml.ID)
						return
					}

					respTml := models.Tml{
						ID:   tml.ID,
						Data: b[:read],
					}

					if err := conn.WriteMessage(websocket.BinaryMessage, respTml.ToBytes()); err != nil {
						if *debug {
							log.Println(err)
						}
						return
					}
				}
			}()

			continue
		//case tml.Close:
		//	p.clientMap[tml.ID].Close()
		//	delete(p.clientMap, tml.ID)
		//	//TODO: 加代码
		//	continue
		default:
			fmt.Println("def: ", tml.ID)
			tcpConn, ex := p.clientMap[tml.ID]
			if !ex {
				// TODO: 改
				fmt.Println("moz")
				return
			}

			_, err := tcpConn.Write(tml.Data)
			if err != nil {
				if *debug {
					log.Println(err)
				}
				// TODO: 改
				p.close(conn, tml.ID)
				fmt.Println("moz  over")
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

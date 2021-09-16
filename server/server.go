package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/dollarkillerx/urllib"
	"github.com/gorilla/websocket"
)

var url string
var target string

func init() {
	url = "piledriver"
	target = "https://www.dollarkiller.com/light/view/"
}

type PiledriverHandler struct{}

//func middlewareAuth(route Handler) Handler {
//	return func (write http.ResponseWriter, req *http.Request) {
//		token := req.Header.Get("token")
//		// 当用户校验通过设置用户的value
//		if token == "token" {
//			//r.ParseForm()
//			//r.PostForm.Set("userid", "dollarkiller")
//			route(write,req)
//			return
//		}
//		write.WriteHeader(401)
//		write.Write([]byte("401"))
//	}
//}

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
	switch req.URL.String() {
	case fmt.Sprintf("/%s", url):
		conn, err := upgrader.Upgrade(write, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Println(data)

			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err)
				return
			}
		}

	default:
		targetUr := fmt.Sprintf("%s/%s", target, req.URL.String())
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
	ser := &PiledriverHandler{}
	err := http.ListenAndServeTLS(":8020", "./key/server.crt", "./key/server.key", ser)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

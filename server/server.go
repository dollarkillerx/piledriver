package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dollarkillerx/urllib"
)

func proxy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Secure Hello World.\n"))
}

var url string
var target string

func init() {
	url = "piledriver"
	target = "https://www.dollarkiller.com/light/view/"
}

type PiledriverHandler struct{}

func (p *PiledriverHandler) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	switch req.URL.String() {
	case fmt.Sprintf("/%s", url):
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

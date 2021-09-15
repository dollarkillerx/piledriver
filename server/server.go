package main

import (
	"fmt"
	"github.com/dollarkillerx/urllib"
	"log"
	"net/http"
)

func proxy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Secure Hello World.\n"))
}

var url string
var target string

func init() {
	url = "piledriver"
	target = "http://www.bequ6.com"
}

type PiledriverHandler struct{}

func (p *PiledriverHandler) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	write.Write([]byte(req.URL.String()))

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
		write.Header().Add("content-type", "text/html charset=utf-8")
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

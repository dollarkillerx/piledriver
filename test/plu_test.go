package test

import (
	"github.com/dollarkillerx/easy_dns"
	"log"
	"net"
	"plumber/storage"
	"strings"
	"testing"
)

//
//func TestIP(t *testing.T) {
//	search, err := utils.IP2.MemorySearch("127.0.0.1")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Println(search.Country)
//	if search.Country == "0" {
//		log.Println("oK")
//	}
//}

func TestDns(t *testing.T) {
	ip, err := net.LookupIP("www.baidu.com")
	if err != nil {
		log.Fatalln(err)
	}
	if len(ip) >= 1 {
		log.Println(ip[0])
	}

}

func TestDns2(t *testing.T) {
	ip, err := easy_dns.LookupIP("www.baidu.com", "8.8.8.8:53")
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range ip.Answers {
		log.Printf("Addr: %s TTL: %d \n", v.Body.GoString(), v.Header.TTL)
	}
}

func TestD2(t *testing.T) {
	c := "127.0.0.1       localhost"
	split := strings.Split(c, " ")
	log.Println(len(split))
}

func TestPac(t *testing.T) {
	get, err := storage.Storage.Get("abc")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(get)
	}

	err = storage.Storage.Set("ada", "adsd")
	if err != nil {
		log.Fatalln(err)
	}

	i, err := storage.Storage.Get("ada")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(i.(string))
}

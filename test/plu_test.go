package test

import (
	"log"
	"net"
	"plumber/utils"
	"testing"
)

func TestIP(t *testing.T) {
	search, err := utils.IP2.MemorySearch("14.215.177.38")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(search.Country)
}

func TestDns(t *testing.T) {
	ip, err := net.LookupIP("www.baidu.com")
	if err != nil {
		log.Fatalln(err)
	}
	if len(ip) >= 1 {
		log.Println(ip[0])
	}
}

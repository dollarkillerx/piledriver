package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"plumber/utils"
	"testing"

	"github.com/dollarkillerx/easy_dns"
)

func TestDns(t *testing.T) {
	ip, err := easy_dns.LookupIP("www.google.com", "8.8.8.8:53")
	if err != nil {
		log.Fatalln(err)
	}

	marshal, err := json.Marshal(ip)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(marshal))
}

func TestDns2(t *testing.T) {
	ip, err := lockDns("www.google.com", "8.8.8.8:53")
	if err != nil {
		log.Fatalln(err)
	}

	marshal, err := json.Marshal(ip)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(marshal))
}

func TestDns3(t *testing.T) {
	host, err := net.LookupHost("chat.mvalley.com")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(host)

	search, err := utils.IP2.MemorySearch(host[0])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(search.Country)
}

func TestDns4(t *testing.T) {
	host, err := lockDns("www.google.com", "8.8.8.8:53")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(host)

	search, err := utils.IP2.MemorySearch(host[0])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(search.Country)

	{
		host, err := lockDns("www.baidu.com", "8.8.8.8:53")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(host)

		search, err := utils.IP2.MemorySearch(host[0])
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(search.Country)

		if search.Country == "中国" {
			fmt.Println(search.Country, "11")
		}
	}
}

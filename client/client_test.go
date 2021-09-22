package main

import (
	"encoding/json"
	"fmt"
	"log"
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

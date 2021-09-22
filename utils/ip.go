package utils

import (
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"

	"log"
)

var IP2 *ip2region.Ip2Region

func init() {
	region, err := ip2region.New("../ip2region.db")
	if err != nil {
		log.Fatalln(err)
	}
	IP2 = region
	log.Println("Ip2Region Success")
}

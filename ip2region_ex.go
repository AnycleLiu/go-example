package main

import (
	"fmt"

	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
)

func main() {
	fmt.Println("err")
	region, err := ip2region.New("data/ip2region.db")
	defer region.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	ip, err := region.MemorySearch("208.27.187.255")
	fmt.Println(ip, err)
	ip, err = region.BinarySearch("185.16.30.3")
	fmt.Println(ip, err)
	ip, err = region.BtreeSearch("31.40.178.0")
	fmt.Println(ip, err)

	ip, err = region.BinarySearch("43.247.60.0")
	fmt.Println(ip, err)

	ip, err = region.BinarySearch("182.239.127.137")
	fmt.Println(ip, err)

	ip, err = region.BinarySearch("219.133.168.5")
	fmt.Println(ip, err)

	ip, err = region.BinarySearch("134.208.0.0")
	fmt.Println(ip, err)
	ip, err = region.BinarySearch("localhost")
	fmt.Println(ip, err)

	ip, err = region.BinarySearch("103.136.251.7")
	fmt.Println(ip, err)

	ip, err = region.BinarySearch("192.168.70.131")
	fmt.Println(ip, err)

}

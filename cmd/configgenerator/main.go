package main

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/datahandling"
)

func main() {
	fmt.Println("Will be implemented soon")

	ipv4, err := datahandling.GetIPv4Address("192.168.100.12/24")
	if err != nil {
		fmt.Println(err)
		return
	}
	ipv4.Calculate()
	fmt.Println(ipv4.GetNetworkAddress())

}

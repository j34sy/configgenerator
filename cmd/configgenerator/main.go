package main

import (
	"fmt"
	"os"

	"github.com/j34sy/configgenerator/pkg/datahandling"
	"github.com/j34sy/configgenerator/pkg/importer"
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

	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path")
		return
	}

	if len(os.Args) > 2 {
		fmt.Println("Please provide only one file path")
		return
	}

	filePath := os.Args[1]

	networks, err := importer.Load(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, network := range *networks {

		importer.PrintFullNetworkInfo(&network)
	}

}

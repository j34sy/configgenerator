package main

import (
	"fmt"
	"os"

	"github.com/j34sy/configgenerator/pkg/devices"
	"github.com/j34sy/configgenerator/pkg/importer"
)

func main() {
	// Check if the user provided a file path
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path")
		return
	}

	if len(os.Args) > 2 {
		fmt.Println("Please provide only one file path")
		return
	}

	filePath := os.Args[1]

	networksYAML, err := importer.LoadYAML(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	networks := []Network{}

	for _, network := range *networksYAML {

		routers := []*devices.Router{}
		switches := []*devices.Switch{}
		mlSwitches := []*devices.MLSwitch{}

		for _, routerYAML := range network.Routers {
			routers = append(routers, devices.CreateRouter(routerYAML, network.Users, network.Name, networksYAML))
		}

		for _, switchYAML := range network.Switches {
			switches = append(switches, devices.CreateSwitch(switchYAML, network.Users, network.Vlans, network.Name))
		}

		for _, mlSwitchYAML := range network.MLSwitches {
			mlSwitches = append(mlSwitches, devices.CreateMLSwitch(mlSwitchYAML, network.Users, networksYAML, network.Name, network.Vlans))
		}

		networkData := Network{
			Routers:    routers,
			Switches:   switches,
			MLSwitches: mlSwitches,
		}

		networks = append(networks, networkData)

	}

	for _, network := range networks {
		fmt.Println("Network: ", network.Routers[0].Domain)
		PrintNetwork(network)
	}

}

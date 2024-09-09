package main

import (
	"fmt"
	"os"

	"github.com/j34sy/configgenerator/cmd/configgenerator/writer"
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

	networks := []writer.Network{}

	for _, network := range *networksYAML {

		routers := []*devices.Router{}
		switches := []*devices.Switch{}
		mlSwitches := []*devices.MLSwitch{}

		for _, routerYAML := range network.Routers {
			routers = append(routers, devices.CreateRouter(routerYAML, network.Users, network.Name, network.EnableSecret, networksYAML))
		}

		for _, switchYAML := range network.Switches {
			switches = append(switches, devices.CreateSwitch(switchYAML, network.Users, network.Vlans, network.Name, network.EnableSecret))
		}

		for _, mlSwitchYAML := range network.MLSwitches {
			mlSwitches = append(mlSwitches, devices.CreateMLSwitch(mlSwitchYAML, network.Users, networksYAML, network.Name, network.EnableSecret, network.Vlans))
		}

		networkData := writer.Network{
			Routers:    routers,
			Switches:   switches,
			MLSwitches: mlSwitches,
			Name:       network.Name,
		}

		networks = append(networks, networkData)

	}

	for _, network := range networks {
		writer.PrintNetwork(network)
		writer.WriteConfigs(network)
	}

}

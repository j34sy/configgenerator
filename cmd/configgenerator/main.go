package main

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/devices"
	"github.com/j34sy/configgenerator/pkg/importer"
)

func main() {
	/*
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

		networksYAML, err := importer.LoadYAML(filePath)

		if err != nil {
			fmt.Println(err)
			return
		}


			for _, network := range *networksYAML {

				importer.PrintFullNetworkInfoByYAML(&network)
			}
	*/

	// Example usage

	fullNetwork := []importer.NetworkYAML{
		{
			Routers: []importer.RouterYAML{
				{
					Name: "Router1",
					Routes: importer.RoutesYAML{
						Default:      "192.168.1.2",
						Destinations: []string{"192.168.2.0/24", "192.168.3.0/24", "192.168.4.0/24", "192.168.5.0/24", "172.16.0.0/24", "192.168.100.0/24", "192.168.200.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.1.1/24"},
						{Name: "eth1", IP: "10.0.0.0/8"},
					},
				},
				{
					Name: "Router2",
					Routes: importer.RoutesYAML{
						Default:      "192.168.2.2",
						Destinations: []string{"192.168.3.0/24", "192.168.4.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24", "192.168.100.0/24", "192.168.200.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.2.1/24"},
						{Name: "eth1", IP: "192.168.1.2/24"},
					},
				},
				{
					Name: "Router3",
					Routes: importer.RoutesYAML{
						Default:      "192.168.3.2",
						Destinations: []string{"192.168.1.0/24", "192.168.4.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24", "192.168.100.0/24", "192.168.200.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.3.1/24"},
						{Name: "eth1", IP: "192.168.2.2/24"},
					},
				},
				{
					Name: "Router4",
					Routes: importer.RoutesYAML{
						Default:      "192.168.4.2",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24", "192.168.100.0/24", "192.168.200.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.4.1/24"},
						{Name: "eth1", IP: "192.168.3.2/24"},
					},
				},
			},
			MLSwitches: []importer.MLSwitchYAML{
				{
					Name: "MLSwitch1",
					Routes: importer.RoutesYAML{
						Default:      "192.168.5.1",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "10.0.0.0/8", "172.16.0.0/24", "192.168.100.0/24", "192.168.200.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.4.2/24"},
						{Name: "eth1", IP: "192.168.5.1/24"},
					},
					Routing: true,
				},
				{
					Name: "MLSwitch2",
					Routes: importer.RoutesYAML{
						Default:      "",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "192.168.4.0/24", "10.0.0.0/8", "192.168.100.0/24", "192.168.200.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.5.2/24"},
						{Name: "eth1", IP: "172.16.0.0/24"},
					},
					Routing: true,
				},
				{
					Name: "MLSwitch3",
					Routes: importer.RoutesYAML{
						Default:      "192.168.4.2",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24", "192.168.200.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.4.3/24"},
						{Name: "eth1", IP: "192.168.100.1/24"},
					},
					Routing: true,
				},
				{
					Name: "MLSwitch4",
					Routes: importer.RoutesYAML{
						Default:      "192.168.100.1",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "192.168.4.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.200.1/24"},
						{Name: "eth1", IP: "192.168.100.2/24"},
					},
					Routing: true,
				},
			},
		},
	}

	dest := "192.168.200.0/24"
	nextHop := devices.FindNextHop(dest, devices.RoutingDevice{
		Name: "Router3",
		Interfaces: []devices.Interface{
			{Name: "eth0", IP: "192.168.3.1/24"},
			{Name: "eth1", IP: "192.168.2.2/24"},
		},
		Destinations: []string{""},
		Default:      "",
	}, &fullNetwork)
	fmt.Println("Next Hop:", nextHop)
}

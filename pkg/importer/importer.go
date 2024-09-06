package importer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func Load(filePath string) (*[]Network, error) {
	// get data from provided file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var networks []Network

	err = yaml.Unmarshal(content, &networks)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &networks, nil
}

func PrintFullNetworkInfo(network *Network) error {
	fmt.Println("Network info for network: " + network.Name)

	fmt.Println()
	fmt.Println("Users:")
	for _, user := range network.Users {
		fmt.Println("User: " + user.Name)
		fmt.Println("Privilege: " + user.Privilege)
		fmt.Println("Password: " + user.Password)
		fmt.Println()
	}
	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Vlans:")
	for _, vlangroup := range network.Vlans {
		fmt.Println("Switches for this VLAN set:")
		for _, switchName := range vlangroup.Switches {
			fmt.Println("Switch: " + switchName)
		}
		fmt.Println("VLANs:")
		for _, vlan := range vlangroup.List {
			fmt.Println("ID: ", vlan.ID)
			fmt.Println("Name: " + vlan.Name)
			fmt.Println("Subnet: " + vlan.Subnet)
			fmt.Println("Gateway: " + vlan.Gateway)
			fmt.Println()
		}
		fmt.Println("---")
	}

	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Switches:")
	for _, switchDev := range network.Switches {
		fmt.Println("Switch: " + switchDev.Name)
		fmt.Println("Interfaces:")
		for _, iface := range switchDev.Interfaces {
			fmt.Println("Name: " + iface.Name)
			fmt.Println("Vlan: " + iface.Vlan)
			fmt.Println("IP: " + iface.IP)
			fmt.Println("Trunk: ")
			for _, trunk := range iface.Trunk {
				fmt.Println(trunk)
			}
			fmt.Println("Access: ", iface.Access)
			fmt.Println("OSPF:")
			if iface.OSPF != nil {
				fmt.Println("Process: ", iface.OSPF.Process)
				fmt.Println("Area: ", iface.OSPF.Area)
			}
			fmt.Println()
		}
		fmt.Println("---")
	}

	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Multi-Layer Switches:")
	for _, mlswitch := range network.MLSwitches {
		fmt.Println("Switch: " + mlswitch.Name)
		fmt.Println("Routing: ", mlswitch.Routing)
		fmt.Println("Interfaces:")
		for _, iface := range mlswitch.Interfaces {
			fmt.Println("Name: " + iface.Name)
			fmt.Println("Vlan: " + iface.Vlan)
			fmt.Println("IP: " + iface.IP)
			fmt.Println("Trunk: ")
			for _, trunk := range iface.Trunk {
				fmt.Println(trunk)
			}
			fmt.Println("Access: ", iface.Access)
			fmt.Println("OSPF:")
			if iface.OSPF != nil {
				fmt.Println("Process: ", iface.OSPF.Process)
				fmt.Println("Area: ", iface.OSPF.Area)
			}
			fmt.Println()
		}
		fmt.Println("---")
	}

	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Routers:")
	for _, router := range network.Routers {
		fmt.Println("Router: " + router.Name)
		fmt.Println("OSPF settings:")
		for _, ospf := range router.OSPFRouter {
			fmt.Println("Process: ", ospf.Process)
			fmt.Println("ID: " + ospf.ID)
			fmt.Println()
		}
		fmt.Println("Interfaces:")
		for _, iface := range router.Interfaces {
			fmt.Println("Name: " + iface.Name)
			fmt.Println("Vlan: " + iface.Vlan)
			fmt.Println("IP: " + iface.IP)
			fmt.Println("Trunk: ")
			for _, trunk := range iface.Trunk {
				fmt.Println(trunk)
			}
			fmt.Println("Access: ", iface.Access)
			fmt.Println("OSPF:")
			if iface.OSPF != nil {
				fmt.Println("Process: ", iface.OSPF.Process)
				fmt.Println("Area: ", iface.OSPF.Area)
			}
			fmt.Println()
		}
		fmt.Println("Routes:")
		for _, route := range router.Routes {
			fmt.Println("Destinations:")
			for _, dest := range route.Destinations {
				fmt.Println(dest)
			}
			fmt.Println("Default: " + route.Default)
			fmt.Println()
		}
		fmt.Println("---")
	}

	fmt.Println("This should be all ...")

	return nil
}

package importer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYAML(filePath string) (*[]NetworkYAML, error) {
	// get data from provided file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var networks []NetworkYAML

	err = yaml.Unmarshal(content, &networks)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &networks, nil
}

func PrintFullNetworkInfoByYAML(network *NetworkYAML) error {
	fmt.Println("Network info for network: " + network.Name)

	fmt.Println()
	fmt.Println("Users:")
	if network.Users != nil {
		for _, user := range network.Users {
			fmt.Println("User: " + user.Name)
			fmt.Println("Privilege: " + user.Privilege)
			fmt.Println("Password: " + user.Password)
			fmt.Println()
		}
	}
	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Vlans:")
	if network.Vlans != nil {
		for _, vlangroup := range network.Vlans {
			fmt.Println("Switches for this VLAN set:")
			for _, switchName := range vlangroup.Switches {
				fmt.Println("Switch: " + switchName)
			}
			fmt.Println("VLANs:")
			if vlangroup.List != nil {
				for _, vlan := range vlangroup.List {
					fmt.Println("ID: ", vlan.ID)
					fmt.Println("Name: " + vlan.Name)
					fmt.Println("Subnet: " + vlan.Subnet)
					fmt.Println("Gateway: " + vlan.Gateway)
					fmt.Println()
				}
			}
			fmt.Println("---")
		}
	}

	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Switches:")
	if network.Switches != nil {
		for _, switchDev := range network.Switches {
			fmt.Println("Switch: " + switchDev.Name)

			if switchDev.Interfaces != nil {

				PrintInterfaces(switchDev.Interfaces)
			}

			fmt.Println("---")
		}
	}

	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Multi-Layer Switches:")
	if network.MLSwitches != nil {

		for _, mlswitch := range network.MLSwitches {
			fmt.Println("Switch: " + mlswitch.Name)
			fmt.Println("Routing: ", mlswitch.Routing)

			if mlswitch.Interfaces != nil {

				PrintInterfaces(mlswitch.Interfaces)
			}

			PrintRoutes(mlswitch.Routes)

			if mlswitch.OSPFRouter != nil {

				PrintOSPFRouter(mlswitch.OSPFRouter)
			}

			fmt.Println("---")

		}
	}

	fmt.Println(". . . . . . . . . . . . . . . . . . . . . .")
	fmt.Println()

	fmt.Println("Routers:")
	for _, router := range network.Routers {
		fmt.Println("Router: " + router.Name)

		if router.OSPFRouter != nil {
			PrintOSPFRouter(router.OSPFRouter)
		}

		if router.Interfaces != nil {

			PrintInterfaces(router.Interfaces)

		}

		PrintRoutes(router.Routes)

		fmt.Println("---")
	}

	fmt.Println("This should be all ...")

	return nil
}

func PrintRoutes(routes RoutesYAML) {
	fmt.Println("Routes:")
	fmt.Println("Destinations:")
	if routes.Destinations != nil {
		for _, dest := range routes.Destinations {
			fmt.Println(dest)
		}
	}
	fmt.Println("Default: " + routes.Default)
	fmt.Println()
}

func PrintInterfaces(interfaces []InterfaceYAML) {
	fmt.Println("Interfaces:")
	for _, iface := range interfaces {
		fmt.Println("Name: " + iface.Name)
		fmt.Println("Vlan: " + iface.Vlan)
		fmt.Println("IP: " + iface.IP)
		fmt.Println("Trunk: ")
		for _, trunk := range iface.Trunk {
			fmt.Println(trunk)
		}
		fmt.Println("Native: ", iface.Native)
		fmt.Println("Access: ", iface.Access)
		fmt.Println("OSPF:")
		if iface.OSPF != nil {
			fmt.Println("Process: ", iface.OSPF.Process)
			fmt.Println("Area: ", iface.OSPF.Area)
		}
		fmt.Println()
	}
}

func PrintOSPFRouter(ospfRouter []OSPFRouterYAML) {
	fmt.Println("OSPF Routers:")
	for _, router := range ospfRouter {
		fmt.Println("Process: ", router.Process)
		fmt.Println("ID: " + router.ID)
		fmt.Println()
	}
}

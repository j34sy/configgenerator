package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/j34sy/configgenerator/pkg/devices"
)

type Network struct {
	Routers    []*devices.Router
	Switches   []*devices.Switch
	MLSwitches []*devices.MLSwitch
	Name       string
}

func PrintNetwork(network Network) {
	fmt.Println("Network:", network.Name)

	fmt.Println("Switches:")
	for _, switchDev := range network.Switches {
		fmt.Println(switchDev.Name, switchDev.Domain)
		fmt.Println("Interfaces:")
		for _, iface := range switchDev.Interfaces {
			fmt.Println(iface.Name, iface.Vlan, iface.IP, iface.Access, iface.Trunk, iface.Native)
		}
		fmt.Println("Users:")
		for _, user := range switchDev.Users {
			fmt.Println(user.Name, user.Password, user.Privilege)
		}
		fmt.Println("Vlans:")
		for _, vlan := range switchDev.Vlans {
			fmt.Println(vlan.ID, vlan.Name)
		}
		fmt.Println("Default Gateway:", switchDev.Default)
	}

	fmt.Println("Routers:")
	for _, router := range network.Routers {
		fmt.Println(router.Name, router.Domain)
		fmt.Println("Interfaces:", len(router.Interfaces))
		for _, iface := range router.Interfaces {
			fmt.Println(iface.Name, iface.Vlan, iface.IP, iface.Access, iface.Trunk, iface.Native)
			if iface.OSPF != nil {
				fmt.Println("OSPF Process:", iface.OSPF.Process)
				fmt.Println("OSPF Area:", iface.OSPF.Area)
			}
		}
		fmt.Println("OSPF Routers:")
		for _, ospfRouter := range router.OSPFRouters {
			fmt.Println(ospfRouter.Process, ospfRouter.ID)
		}

		fmt.Println("Users:")
		for _, user := range router.Users {
			fmt.Println(user.Name, user.Password, user.Privilege)
		}

		fmt.Println("Routes:")
		for _, route := range router.Routes {
			fmt.Println(route.Destination, route.NextHop)
		}

		fmt.Println("Default Gateway:", router.Default)
	}

	fmt.Println("MLSwitches:")
	for _, mlSwitch := range network.MLSwitches {
		fmt.Println(mlSwitch.Name, mlSwitch.Domain)
		fmt.Println("Interfaces:")
		for _, iface := range mlSwitch.Interfaces {
			fmt.Println(iface.Name, iface.Vlan, iface.IP, iface.Access, iface.Trunk, iface.Native)
			if iface.OSPF != nil {
				fmt.Println("OSPF Process:", iface.OSPF.Process)
				fmt.Println("OSPF Area:", iface.OSPF.Area)
			}
		}

		fmt.Println("Routing:", mlSwitch.Routing)

		fmt.Println("OSPF Routers:")
		for _, ospfRouter := range mlSwitch.OSPFRouters {
			fmt.Println(ospfRouter.Process, ospfRouter.ID)
		}

		fmt.Println("Users:")
		for _, user := range mlSwitch.Users {
			fmt.Println(user.Name, user.Password, user.Privilege)
		}

		fmt.Println("Routes:")
		for _, route := range mlSwitch.Routes {
			fmt.Println(route.Destination, route.NextHop)
		}

		fmt.Println("Vlans:")
		for _, vlan := range mlSwitch.Vlans {
			fmt.Println(vlan.ID, vlan.Name)
		}

		fmt.Println("Default Gateway:", mlSwitch.Default)
	}
}

func WriteConfigs(network Network) {
	dirName := "configs_" + network.Name

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			fmt.Println("Error creating directory: ", err)
			return
		}
	}

	dirPath, err := filepath.Abs(dirName)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	for _, router := range network.Routers {
		writeRouterConfigFile(router, dirPath)
	}

	for _, switchDev := range network.Switches {
		WriteSwitchConfigFile(switchDev, dirPath)
	}

	for _, mlSwitch := range network.MLSwitches {
		WriteMLSwitchConfigFile(mlSwitch, dirPath)
	}

	fmt.Println("Configs written to:", dirPath)

}

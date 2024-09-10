package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/j34sy/configgenerator/pkg/devices"
)

func writeRouterConfigFile(router *devices.Router, dirPath string) {
	fileName := router.Name + "." + router.Domain + ".cisco"

	filePath := filepath.Join(dirPath, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating or clearing file:", err)
		return
	}
	file.Close()

	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	file.WriteString("hostname " + router.Name + "\n")
	file.WriteString("ip domain-name " + router.Domain + "\n")
	file.WriteString("enable secret " + router.EnableSecret + "\n")
	file.WriteString("ipv6 unicast-routing\n")
	writeUsers(file, router.Users)
	file.WriteString("line console 0\n")
	file.WriteString("logging synchronous\n")
	if len(router.Users) > 0 {
		file.WriteString("login local\n")
	}
	file.WriteString("no ip domain-lookup\n")

	for _, ospfRouter := range router.OSPFRouters {
		file.WriteString("router ospf " + strconv.Itoa(ospfRouter.Process) + "\n")
		file.WriteString("router-id " + ospfRouter.ID + "\n")
		// TODO: Not adding ospf info on interface but here..., needs data remodeling
		file.WriteString("exit\n")
		file.WriteString("ipv6 router ospf " + strconv.Itoa(ospfRouter.Process) + "\n")
		file.WriteString("router-id " + ospfRouter.ID + "\n")
		// TODO: Not adding ospf info on interface but here..., needs data remodeling
		file.WriteString("exit\n")
	}

	subifaces := []string{}

	for _, iface := range router.Interfaces {
		if iface.Vlan != "" {
			subifaces = append(subifaces, iface.Name)
			writeInterfaceLayer3(file, iface, iface.Name+"."+iface.Vlan)
		} else {
			writeInterfaceLayer3(file, iface, iface.Name)
		}
	}

	subifaces = getAdditionalInterfaces(subifaces)
	//fmt.Println("Subinterfaces: ", subifaces, " for router: ", router.Name)

	for _, subiface := range subifaces {
		found := false
		for _, iface := range router.Interfaces {
			//fmt.Println("Checking interface: ", iface.Name)
			if iface.Name == subiface && iface.Vlan == "" {
				found = true

			}
		}
		if !found {
			//fmt.Println("Creating subinterface: ", subiface)
			file.WriteString("interface " + subiface + "\n")
			file.WriteString("no shutdown \n")
			file.WriteString("exit \n")
		}

	}

	writeRoutesv4(file, router.Routes)
	writeRoutesv6(file, router.Routesv6)

	if router.Default != "" {
		file.WriteString("ip route 0.0.0.0 0.0.0.0 " + strings.Split(router.Default, "/")[0] + "\n")
	}

	if router.Defaultv6 != "" {
		file.WriteString("ipv6 route ::/0 " + strings.Split(router.Defaultv6, "/")[0] + "\n")
	}

	writeSSH(file)

}

func getAdditionalInterfaces(subifaces []string) []string {
	shortedSubIfaces := []string{}

	for _, iface := range subifaces {
		if !containsElement(shortedSubIfaces, iface) {
			shortedSubIfaces = append(shortedSubIfaces, iface)
		}
	}
	return shortedSubIfaces

}

func containsElement(list []string, element string) bool {
	for _, elem := range list {
		if elem == element {
			return true
		}
	}
	return false
}

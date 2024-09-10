package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/j34sy/configgenerator/pkg/devices"
)

func WriteMLSwitchConfigFile(mlSwitch *devices.MLSwitch, dirPath string) {
	fileName := mlSwitch.Name + "." + mlSwitch.Domain + ".cisco"

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

	file.WriteString("hostname " + mlSwitch.Name + "\n")
	file.WriteString("ip domain-name " + mlSwitch.Domain + "\n")
	file.WriteString("enable secret " + mlSwitch.EnableSecret + "\n")
	writeUsers(file, mlSwitch.Users)
	file.WriteString("line console 0\n")
	file.WriteString("logging synchronous\n")
	if len(mlSwitch.Users) > 0 {
		file.WriteString("login local\n")
	}

	file.WriteString("no ip domain-lookup\n")

	writeVlans(file, mlSwitch.Vlans)

	if mlSwitch.Routing {

		// TODO: setup routing interfaces...
		file.WriteString("ip routing\n")
		file.WriteString("ipv6 unicast-routing\n")

		for _, ospfRouter := range mlSwitch.OSPFRouters {
			file.WriteString("router ospf " + strconv.Itoa(ospfRouter.Process) + "\n")
			file.WriteString("router-id " + ospfRouter.ID + "\n")
			// TODO: Not adding ospf info on interface but here..., needs data remodeling
			file.WriteString("exit\n")
			file.WriteString("ipv6 router ospf " + strconv.Itoa(ospfRouter.Process) + "\n")
			file.WriteString("router-id " + ospfRouter.ID + "\n")
			// TODO: Not adding ospf info on interface but here..., needs data remodeling
			file.WriteString("exit\n")
		}

		for _, iface := range mlSwitch.Interfaces {
			writeInterfaceLayer3(file, iface, iface.Name)
		}

		writeRoutesv4(file, mlSwitch.Routes)
		writeRoutesv6(file, mlSwitch.Routesv6)

		if mlSwitch.Default != "" {
			file.WriteString("ip route 0.0.0.0 0.0.0.0 " + strings.Split(mlSwitch.Default, "/")[0] + "\n")
		}

		if mlSwitch.Defaultv6 != "" {
			file.WriteString("ipv6 route ::/0 " + strings.Split(mlSwitch.Defaultv6, "/")[0] + "\n")
		}

	} else {
		for _, iface := range mlSwitch.Interfaces {
			writeInterfaceLayer2(file, iface)
		}

		if mlSwitch.Default != "" {
			file.WriteString("ip default-gateway " + splitIP(mlSwitch.Default) + "\n")
		}

		if mlSwitch.Defaultv6 != "" {
			// FIXME: Find correct way of setting default gw for ipv6
			// folowing docs
			// file.WriteString("ipv6 default-gateway " + switchDev.Defaultv6 + "\n")
			file.WriteString("ipv6 route ::/0 " + splitIP(mlSwitch.Defaultv6) + "\n")
		}
	}

	writeSSH(file)

}

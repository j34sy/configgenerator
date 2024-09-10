package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/j34sy/configgenerator/pkg/devices"
)

func WriteSwitchConfigFile(switchDev *devices.Switch, dirPath string) {
	fileName := switchDev.Name + "." + switchDev.Domain + ".cisco"

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

	file.WriteString("hostname " + switchDev.Name + "\n")
	file.WriteString("ip domain-name " + switchDev.Domain + "\n")
	file.WriteString("enable secret " + switchDev.EnableSecret + "\n")
	writeUsers(file, switchDev.Users)
	file.WriteString("line console 0\n")
	file.WriteString("logging synchronous\n")
	if len(switchDev.Users) > 0 {
		file.WriteString("login local\n")
	}

	file.WriteString("no ip domain-lookup\n")

	writeVlans(file, switchDev.Vlans)

	for _, iface := range switchDev.Interfaces {
		writeInterfaceLayer2(file, iface)
	}

	if switchDev.Default != "" {
		file.WriteString("ip default-gateway " + splitIP(switchDev.Default) + "\n")
	}

	if switchDev.Defaultv6 != "" {
		// FIXME: Find correct way of setting default gw for ipv6
		// folowing docs
		// file.WriteString("ipv6 default-gateway " + switchDev.Defaultv6 + "\n")
		file.WriteString("ipv6 route ::/0 " + splitIP(switchDev.Defaultv6) + "\n")
	}

	writeSSH(file)
}

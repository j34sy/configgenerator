package writer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/j34sy/configgenerator/pkg/datahandling"
	"github.com/j34sy/configgenerator/pkg/devices"
)

func writeInterfaceLayer2(file *os.File, iface devices.Interface) {
	file.WriteString("interface " + iface.Name + "\n")
	if iface.Vlan == "access" {
		file.WriteString("switchport mode access\n")
		file.WriteString("switchport access vlan " + strconv.Itoa(iface.Access) + "\n")
	} else if iface.Vlan == "trunk" {
		file.WriteString("switchport mode trunk\n")
		file.WriteString("switchport trunk native vlan " + strconv.Itoa(iface.Native) + "\n")
		file.WriteString("switchport trunk allowed vlan " + strconv.Itoa(iface.Native) + "\n")
		for _, trunk := range iface.Trunk {
			file.WriteString("switchport trunk allowed vlan add " + strconv.Itoa(trunk) + "\n")
		}
	}

	checkVlanInt := devices.IsVlanInterface(iface.Name)
	if checkVlanInt {
		_, err := devices.ExtractVLAN(iface.Name)
		if err != nil {
			fmt.Println("Error extracting vlan ID: ", err)
			return
		}

		if iface.IP != "" {
			devIP, err := datahandling.GetIPv4Address(iface.IP)
			if err != nil {
				fmt.Println("Error getting IPv4 address: ", err)
				return
			}

			subnet := devIP.GetSubnetMask()
			ip := splitIP(iface.IP)

			file.WriteString("ip address " + ip + " " + subnet + "\n")
		}

		if iface.IPv6 != "" {
			file.WriteString("ipv6 address " + iface.IPv6 + "\n")
		}
	}

	file.WriteString("no shutdown\n")
	file.WriteString("exit\n")
}

func splitIP(ipWithCIDR string) string {
	ip := strings.Split(ipWithCIDR, "/")[0]
	return ip
}

func writeInterfaceLayer3(file *os.File, iface devices.Interface, interfaceName string) {
	file.WriteString("interface " + interfaceName + "\n")

	if iface.IP != "" {
		devIP, err := datahandling.GetIPv4Address(iface.IP)
		if err != nil {
			fmt.Println("Error getting IPv4 address: ", err)
			return
		}

		subnet := devIP.GetSubnetMask()
		ip := splitIP(iface.IP)

		file.WriteString("ip address " + ip + " " + subnet + "\n")
	}

	if iface.IPv6 != "" {
		file.WriteString("ipv6 address " + iface.IPv6 + "\n")
	}

	file.WriteString("no shutdown\n")

	if iface.OSPF != nil {
		file.WriteString("ospf " + strconv.Itoa(iface.OSPF.Process) + " area " + strconv.Itoa(iface.OSPF.Area) + "\n")
	}

	if iface.Name != interfaceName && strings.Contains(interfaceName, ".") {
		if iface.Vlan == strconv.Itoa(iface.Native) {
			file.WriteString("encapsulation dot1Q " + iface.Vlan + " native \n")
		} else {
			file.WriteString("encapsulation dot1Q " + iface.Vlan + "\n")

		}
	}

	file.WriteString("exit\n")
}

func writeVlans(file *os.File, vlans []devices.Vlan) {
	for _, vlan := range vlans {
		file.WriteString("vlan " + strconv.Itoa(vlan.ID) + "\n")
		file.WriteString("name " + vlan.Name + "\n")
		file.WriteString("exit\n")
	}
}

func writeUsers(file *os.File, users []devices.User) {
	for _, user := range users {
		file.WriteString("username " + user.Name + " privilege " + user.Privilege + " secret " + user.Password + "\n")
	}
}

func writeSSH(file *os.File) {
	file.WriteString("crypto key generate rsa general-keys modulus 4096\n")
	file.WriteString("ip ssh version 2\n")
	file.WriteString("line vty 0 4\n")
	file.WriteString("transport input ssh\n")
	file.WriteString("login local\n")
	file.WriteString("exit\n")
}

func writeRoutesv4(file *os.File, routes []devices.Route) {
	for _, route := range routes {
		ipv4, err := datahandling.GetIPv4Address(route.Destination)
		if err != nil {
			fmt.Println("Error getting IPv4 address: ", err)
			return
		}

		file.WriteString("ip route " + ipv4.GetNetworkAddress() + " " + ipv4.GetSubnetMask() + " " + route.NextHop + "\n")
	}
}

func writeRoutesv6(file *os.File, routes []devices.Routev6) {
	for _, route := range routes {
		file.WriteString("ipv6 route " + route.Destination + " " + strings.Split(route.NextHop, "/")[0] + "\n")
	}
}

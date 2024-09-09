package devices

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateSwitch(switchYAML importer.SwitchYAML, usersYAML []importer.UserYAML, vlanGroupsYAML []importer.VlanGroupYAML, domain string) *Switch {

	//fmt.Println("Creating switch: ", switchYAML.Name)
	//fmt.Println("in domain: ", domain)

	switchVlans := []Vlan{}
	for _, vlanGroup := range vlanGroupsYAML {
		if contains(vlanGroup.Switches, switchYAML.Name) {
			for _, vlan := range vlanGroup.List {
				switchVlans = append(switchVlans, Vlan(vlan))
			}
		}
	}

	users := []User{}

	for _, userYAML := range usersYAML {
		users = append(users, User(userYAML))
	}

	interfaces := []Interface{}

	for _, iface := range switchYAML.Interfaces {
		var ospf *OSPF
		if iface.OSPF != nil {
			//fmt.Println("Found ospf info in interface: ", iface.Name, " device: ", switchYAML.Name)
			ospf = &OSPF{iface.OSPF.Process, iface.OSPF.Area}

		}
		interfaces = append(interfaces, Interface{iface.Name, iface.Vlan, iface.IP, iface.Trunk, iface.Access, ospf, iface.Native})
	}

	vlanInterfaces := []int{}
	defaultGateway := ""

	for _, iface := range interfaces {
		for _, vlan := range switchVlans {
			if isVlanInterface(iface.Name) {
				vlanID, err := extractVLAN(iface.Name)
				if err != nil {
					fmt.Println("Error extracting vlan ID: ", err)
					continue
				}
				if vlan.ID == vlanID {
					//fmt.Println("Found vlan info in interface: ", iface.Name, " device: ", switchYAML.Name)
					if iface.IP != "" {
						//fmt.Println("Found IP info in interface: ", iface.Name, " device: ", switchYAML.Name)
						vlanInterfaces = append(vlanInterfaces, vlan.ID)
					}
				}
			}
		}
	}

	if len(vlanInterfaces) == 0 {
		//fmt.Println("No vlan interfaces found for device: ", switchYAML.Name)
	} else if len(vlanInterfaces) == 1 {
		//fmt.Println("Found one vlan interface for device: ", switchYAML.Name)
		//fmt.Println("Vlan interface: ", vlanInterfaces[0])
		for _, vlan := range switchVlans {
			if vlan.ID == vlanInterfaces[0] {
				defaultGateway = vlan.Gateway
				//fmt.Println("Default gateway: ", defaultGateway)
			}
		}
	}

	return &Switch{
		Name:       switchYAML.Name,
		Interfaces: interfaces,
		Vlans:      switchVlans,
		Users:      users,
		Default:    defaultGateway,
		Domain:     domain,
	}
}

func isVlanInterface(s string) bool {
	re := regexp.MustCompile(`^vlan\d{1,4}$`)
	return re.MatchString(s)
}

func extractVLAN(s string) (int, error) {
	re := regexp.MustCompile(`^vlan(\d{1,4})$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return strconv.Atoi(matches[1])
	}
	return 0, fmt.Errorf("invalid VLAN string")
}

package devices

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateMLSwitch(mlSwitchYAML importer.MLSwitchYAML, usersYAML []importer.UserYAML, fullNetwork *[]importer.NetworkYAML, domain string, vlanGroupsYAML []importer.VlanGroupYAML) *MLSwitch {

	fmt.Println("Creating MLSwitch: ", mlSwitchYAML.Name)
	fmt.Println("in domain: ", domain)

	interfaces := []Interface{}

	for _, iface := range mlSwitchYAML.Interfaces {
		var ospf *OSPF
		if iface.OSPF != nil {
			ospf = &OSPF{iface.OSPF.Process, iface.OSPF.Area}
			fmt.Println("Found ospf info in interface: ", iface.Name, " device: ", mlSwitchYAML.Name)
		}
		interfaces = append(interfaces, Interface{iface.Name, iface.Vlan, iface.IP, iface.Trunk, iface.Access, ospf, iface.Native})
	}

	ospfRouters := []OSPFRouter{}
	for _, ospfrouter := range mlSwitchYAML.OSPFRouter {
		ospfRouters = append(ospfRouters, OSPFRouter(ospfrouter))
	}

	users := []User{}
	for _, userYAML := range usersYAML {
		users = append(users, User(userYAML))
	}

	destinations := []string{}
	destinations = append(destinations, mlSwitchYAML.Routes.Destinations...)

	routes := []Route{}
	for _, destination := range destinations {
		nextHop := FindNextHop(destination, RoutingDevice{mlSwitchYAML.Name, interfaces, mlSwitchYAML.Routes.Destinations, mlSwitchYAML.Routes.Default}, fullNetwork)
		routes = append(routes, Route{destination, nextHop})
	}

	switchVlans := []Vlan{}
	for _, vlanGroup := range vlanGroupsYAML {
		if contains(vlanGroup.Switches, mlSwitchYAML.Name) {
			for _, vlan := range vlanGroup.List {
				switchVlans = append(switchVlans, Vlan(vlan))
			}
		}
	}

	if !mlSwitchYAML.Routing && mlSwitchYAML.Routes.Default == "" {
		fmt.Println("No routing enabled for device: ", mlSwitchYAML.Name)
		fmt.Println("Looking for default gateway in vlan interfaces")
		vlanInterfaces := []int{}
		defaultGateway := ""

		for _, iface := range interfaces {
			for _, vlan := range switchVlans {
				if iface.Vlan == fmt.Sprint(vlan.ID) {
					fmt.Println("Found vlan info in interface: ", iface.Name, " device: ", mlSwitchYAML.Name)
					if iface.IP != "" {
						fmt.Println("Found IP info in interface: ", iface.Name, " device: ", mlSwitchYAML.Name)
						vlanInterfaces = append(vlanInterfaces, vlan.ID)
					}
				}
			}
		}

		if len(vlanInterfaces) == 0 {
			fmt.Println("No vlan interfaces found for device: ", mlSwitchYAML.Name)
		} else if len(vlanInterfaces) == 1 {
			fmt.Println("Found one vlan interface for device: ", mlSwitchYAML.Name)
			fmt.Println("Vlan interface: ", vlanInterfaces[0])
			for _, vlan := range switchVlans {
				if vlan.ID == vlanInterfaces[0] {
					defaultGateway = vlan.Gateway
					fmt.Println("Default gateway: ", defaultGateway)
				}
			}
		}
	}

	return &MLSwitch{
		RoutingDevice: RoutingDevice{
			mlSwitchYAML.Name,
			interfaces,
			mlSwitchYAML.Routes.Destinations,
			mlSwitchYAML.Routes.Default,
		},
		OSPFRouters: ospfRouters,
		Routing:     mlSwitchYAML.Routing,
		Users:       users,
		Routes:      routes,
		Vlans:       switchVlans,
		Domain:      domain,
	}
}

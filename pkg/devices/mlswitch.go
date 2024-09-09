package devices

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateMLSwitch(mlSwitchYAML importer.MLSwitchYAML, usersYAML []importer.UserYAML, fullNetwork *[]importer.NetworkYAML, domain string, enableSecret string, vlanGroupsYAML []importer.VlanGroupYAML) *MLSwitch {

	//fmt.Println("Creating MLSwitch: ", mlSwitchYAML.Name)
	//fmt.Println("in domain: ", domain)

	interfaces := []Interface{}

	for _, iface := range mlSwitchYAML.Interfaces {
		var ospf *OSPF
		if iface.OSPF != nil {
			ospf = &OSPF{iface.OSPF.Process, iface.OSPF.Area}
			//fmt.Println("Found ospf info in interface: ", iface.Name, " device: ", mlSwitchYAML.Name)
		}
		interfaces = append(interfaces, Interface{iface.Name, iface.Vlan, iface.IP, iface.IPv6, iface.Trunk, iface.Access, ospf, iface.Native})
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
		nextHop := FindNextHop(destination, RoutingDevice{mlSwitchYAML.Name, interfaces, mlSwitchYAML.Routes.Destinations, []string{}, mlSwitchYAML.Routes.Default, ""}, fullNetwork)
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

	var gw string
	var gwv6 string

	if mlSwitchYAML.Routing {
		gw = mlSwitchYAML.Routes.Default
		gwv6 = mlSwitchYAML.Routes.Defaultv6
	}

	routesv6 := []Routev6{}
	destinationsv6 := []string{}
	destinationsv6 = append(destinationsv6, mlSwitchYAML.Routes.Destinationsv6...)

	for _, destination := range destinationsv6 {
		nextHop := FindNextHopv6(destination, RoutingDevice{mlSwitchYAML.Name, interfaces, mlSwitchYAML.Routes.Destinations, mlSwitchYAML.Routes.Destinationsv6, mlSwitchYAML.Routes.Default, mlSwitchYAML.Routes.Defaultv6}, fullNetwork)
		routesv6 = append(routesv6, Routev6{destination, nextHop})
	}

	if !mlSwitchYAML.Routing && mlSwitchYAML.Routes.Default == "" {
		//fmt.Println("No routing enabled for device: ", mlSwitchYAML.Name)
		//fmt.Println("Looking for default gateway in vlan interfaces")
		vlanInterfaces := []int{}
		defaultGateway := ""

		for _, iface := range interfaces {
			for _, vlan := range switchVlans {
				if iface.Vlan == fmt.Sprint(vlan.ID) {
					//fmt.Println("Found vlan info in interface: ", iface.Name, " device: ", mlSwitchYAML.Name)
					if iface.IP != "" {
						//fmt.Println("Found IP info in interface: ", iface.Name, " device: ", mlSwitchYAML.Name)
						vlanInterfaces = append(vlanInterfaces, vlan.ID)
					}
				}
			}
		}

		if len(vlanInterfaces) == 0 {
			//fmt.Println("No vlan interfaces found for device: ", mlSwitchYAML.Name)
		} else if len(vlanInterfaces) == 1 {
			//fmt.Println("Found one vlan interface for device: ", mlSwitchYAML.Name)
			//fmt.Println("Vlan interface: ", vlanInterfaces[0])
			for _, vlan := range switchVlans {
				if vlan.ID == vlanInterfaces[0] {
					defaultGateway = vlan.Gateway
					mlSwitchYAML.Routes.Default = defaultGateway
					gw = defaultGateway
					//fmt.Println("Default gateway: ", defaultGateway)
				}
			}
		}
	}

	if !mlSwitchYAML.Routing && mlSwitchYAML.Routes.Defaultv6 == "" {
		vlanInterfacesv6 := []int{}
		defaultGatewayv6 := ""

		for _, iface := range interfaces {
			for _, vlan := range switchVlans {
				if iface.Vlan == fmt.Sprint(vlan.ID) {
					if iface.IPv6 != "" {
						vlanInterfacesv6 = append(vlanInterfacesv6, vlan.ID)
					}
				}
			}
		}
		if len(vlanInterfacesv6) == 0 {
			//fmt.Println("No vlan interfaces found for device: ", mlSwitchYAML.Name)
		} else if len(vlanInterfacesv6) == 1 {
			for _, vlan := range switchVlans {
				if vlan.ID == vlanInterfacesv6[0] {
					defaultGatewayv6 = vlan.Gatewayv6
					mlSwitchYAML.Routes.Defaultv6 = defaultGatewayv6
					gwv6 = defaultGatewayv6
				}
			}
		}
	}

	return &MLSwitch{
		RoutingDevice: RoutingDevice{
			mlSwitchYAML.Name,
			interfaces,
			mlSwitchYAML.Routes.Destinations,
			mlSwitchYAML.Routes.Destinationsv6,
			gw,
			gwv6,
		},
		OSPFRouters: ospfRouters,
		Routing:     mlSwitchYAML.Routing,
		Users:       users,
		Routes:      routes,
		Routesv6:    routesv6,
		Vlans:       switchVlans,
		Domain:      domain,
	}
}

package devices

import (
	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateRouter(routerYAML importer.RouterYAML, usersYAML []importer.UserYAML, domain string, fullNetwork *[]importer.NetworkYAML) *Router {

	//fmt.Println("Creating router: ", routerYAML.Name)
	//fmt.Println("in domain: ", domain)

	users := []User{}

	for _, userYAML := range usersYAML {
		users = append(users, User(userYAML))
	}

	interfaces := []Interface{}

	for _, iface := range routerYAML.Interfaces {
		var ospf *OSPF
		if iface.OSPF != nil {
			//fmt.Println("Found ospf info in interface: ", iface.Name, " device: ", routerYAML.Name)
			ospf = &OSPF{iface.OSPF.Process, iface.OSPF.Area}

		}
		interfaces = append(interfaces, Interface{iface.Name, iface.Vlan, iface.IP, iface.Trunk, iface.Access, ospf, iface.Native})
	}

	ospfRouters := []OSPFRouter{}

	for _, ospfRouter := range routerYAML.OSPFRouter {
		ospfRouters = append(ospfRouters, OSPFRouter(ospfRouter))
	}

	destinations := []string{}

	destinations = append(destinations, routerYAML.Routes.Destinations...)

	routes := []Route{}

	for _, destination := range destinations {
		nextHop := FindNextHop(destination, RoutingDevice{routerYAML.Name, interfaces, routerYAML.Routes.Destinations, routerYAML.Routes.Default}, fullNetwork)

		routes = append(routes, Route{destination, nextHop})
	}

	return &Router{
		RoutingDevice: RoutingDevice{
			routerYAML.Name,
			interfaces,
			routerYAML.Routes.Destinations,
			routerYAML.Routes.Default,
		},
		Domain:      domain,
		Users:       users,
		OSPFRouters: ospfRouters,
		Routes:      routes,
	}
}

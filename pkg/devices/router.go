package devices

import (
	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateRouter(routerYAML importer.RouterYAML, usersYAML []importer.UserYAML, domain string, enableSecret string, fullNetwork *[]importer.NetworkYAML) *Router {

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
		interfaces = append(interfaces, Interface{iface.Name, iface.Vlan, iface.IP, iface.IPv6, iface.Trunk, iface.Access, ospf, iface.Native})
	}

	ospfRouters := []OSPFRouter{}

	for _, ospfRouter := range routerYAML.OSPFRouter {
		ospfRouters = append(ospfRouters, OSPFRouter(ospfRouter))
	}

	destinations := []string{}

	destinations = append(destinations, routerYAML.Routes.Destinations...)

	routes := []Route{}

	for _, destination := range destinations {
		nextHop := FindNextHop(destination, RoutingDevice{routerYAML.Name, interfaces, routerYAML.Routes.Destinations, []string{}, routerYAML.Routes.Default, ""}, fullNetwork)

		routes = append(routes, Route{destination, nextHop})
	}

	routesv6 := []Routev6{}
	destinationsv6 := []string{}

	destinationsv6 = append(destinationsv6, routerYAML.Routes.Destinationsv6...)

	for _, destination := range destinationsv6 {
		nextHop := FindNextHopv6(destination, RoutingDevice{routerYAML.Name, interfaces, routerYAML.Routes.Destinations, routerYAML.Routes.Destinationsv6, routerYAML.Routes.Default, routerYAML.Routes.Defaultv6}, fullNetwork)
		routesv6 = append(routesv6, Routev6{destination, nextHop})
	}

	return &Router{
		RoutingDevice: RoutingDevice{
			routerYAML.Name,
			interfaces,
			routerYAML.Routes.Destinations,
			routerYAML.Routes.Destinationsv6,
			routerYAML.Routes.Default,
			routerYAML.Routes.Defaultv6,
		},
		Domain:       domain,
		Users:        users,
		OSPFRouters:  ospfRouters,
		Routes:       routes,
		Routesv6:     routesv6,
		EnableSecret: enableSecret,
	}
}

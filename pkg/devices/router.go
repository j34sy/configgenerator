package devices

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/datahandling"
	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateRouter(routerYAML importer.RouterYAML, usersYAML []importer.UserYAML, domain string, fullNetwork *[]importer.NetworkYAML) *Router {

	fmt.Println("Creating router: ", routerYAML.Name)
	fmt.Println("in domain: ", domain)

	users := []User{}

	for _, userYAML := range usersYAML {
		users = append(users, User(userYAML))
	}

	interfaces := []Interface{}

	for _, iface := range routerYAML.Interfaces {
		var ospf *OSPF
		if iface.OSPF != nil {
			fmt.Println("Found ospf info in interface: ", iface.Name, " device: ", routerYAML.Name)
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
		routes = append(routes, Route{destination, ""})
		FindNextHop(destination, routerYAML, fullNetwork)
	}

	return &Router{
		Name:        routerYAML.Name,
		Interfaces:  interfaces,
		Routes:      routes,
		OSPFRouters: ospfRouters,
		Users:       users,
		Default:     routerYAML.Routes.Default,
		Domain:      domain,
	}
}

// Function to get the direct neighbors of each router
func GetDirectNeighbors(router importer.RouterYAML, fullNetwork *[]importer.NetworkYAML) map[string]string {
	neighbors := make(map[string]string)

	for _, network := range *fullNetwork {
		for _, r := range network.Routers {
			if r.Name == router.Name {
				continue
			}
			for _, iface := range r.Interfaces {
				for _, routerIface := range router.Interfaces {
					check, err := datahandling.IsSameNetwork(routerIface.IP, iface.IP)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if check {
						neighbors[r.Name] = iface.IP
					}
				}
			}
		}
		for _, mlswitch := range network.MLSwitches {
			for _, iface := range mlswitch.Interfaces {
				for _, routerIface := range router.Interfaces {
					check, err := datahandling.IsSameNetwork(routerIface.IP, iface.IP)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if check {
						neighbors[mlswitch.Name] = iface.IP
					}
				}
			}
		}
	}

	return neighbors
}

// Function to find the next hop recursively
func FindNextHop(dest string, router importer.RouterYAML, fullNetwork *[]importer.NetworkYAML) string {
	visited := make(map[string]bool)
	return findNextHopRecursive(dest, router, fullNetwork, visited)
}

// Recursive function to find the next hop
func findNextHopRecursive(dest string, router importer.RouterYAML, fullNetwork *[]importer.NetworkYAML, visited map[string]bool) string {
	// Mark the current router as visited
	visited[router.Name] = true

	neighbors := GetDirectNeighbors(router, fullNetwork)

	// Check direct neighbors first
	for neighborName, neighborIP := range neighbors {
		neighborRouter := getRouterByName(neighborName, fullNetwork)
		for _, iface := range neighborRouter.Interfaces {
			check, err := datahandling.IsSameNetwork(iface.IP, dest)
			if err != nil {
				fmt.Println(err)
				return router.Routes.Default
			}
			if check {
				return neighborIP
			}
		}
	}

	// If no direct route is found, check for multi-hop routes
	for neighborName, neighborIP := range neighbors {
		if !visited[neighborName] {
			neighborRouter := getRouterByName(neighborName, fullNetwork)
			nextHop := findNextHopRecursive(dest, neighborRouter, fullNetwork, visited)
			if nextHop != router.Routes.Default {
				return neighborIP
			}
		}
	}

	return router.Routes.Default
}

// Helper function to get the router by its name
func getRouterByName(name string, fullNetwork *[]importer.NetworkYAML) importer.RouterYAML {
	for _, network := range *fullNetwork {
		for _, r := range network.Routers {
			if r.Name == name {
				return r
			}
		}
	}
	return importer.RouterYAML{}
}

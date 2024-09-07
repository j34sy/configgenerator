package devices

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/datahandling"
	"github.com/j34sy/configgenerator/pkg/importer"
)

// Function to get the direct neighbors of each router and MLSwitch
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

// Function to get the direct neighbors of each MLSwitch
func GetDirectNeighborsMLSwitch(mlswitch importer.MLSwitchYAML, fullNetwork *[]importer.NetworkYAML) map[string]string {
	neighbors := make(map[string]string)

	for _, network := range *fullNetwork {
		for _, r := range network.Routers {
			for _, iface := range r.Interfaces {
				for _, mlswitchIface := range mlswitch.Interfaces {
					check, err := datahandling.IsSameNetwork(mlswitchIface.IP, iface.IP)
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
		for _, mls := range network.MLSwitches {
			if mls.Name == mlswitch.Name {
				continue
			}
			for _, iface := range mls.Interfaces {
				for _, mlswitchIface := range mlswitch.Interfaces {
					check, err := datahandling.IsSameNetwork(mlswitchIface.IP, iface.IP)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if check {
						neighbors[mls.Name] = iface.IP
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

func FindNextHopML(dest string, mlswitch importer.MLSwitchYAML, fullNetwork *[]importer.NetworkYAML) string {
	visited := make(map[string]bool)
	return findNextHopRecursiveMLSwitch(dest, mlswitch, fullNetwork, visited)
}

// Recursive function to find the next hop for routers
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

// Recursive function to find the next hop for MLSwitches
func findNextHopRecursiveMLSwitch(dest string, mlswitch importer.MLSwitchYAML, fullNetwork *[]importer.NetworkYAML, visited map[string]bool) string {
	// Mark the current MLSwitch as visited
	visited[mlswitch.Name] = true

	neighbors := GetDirectNeighborsMLSwitch(mlswitch, fullNetwork)

	// Check direct neighbors first
	for neighborName, neighborIP := range neighbors {
		neighborRouter := getRouterByName(neighborName, fullNetwork)
		for _, iface := range neighborRouter.Interfaces {
			check, err := datahandling.IsSameNetwork(iface.IP, dest)
			if err != nil {
				fmt.Println(err)
				return mlswitch.Routes.Default
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
			if nextHop != mlswitch.Routes.Default {
				return neighborIP
			}
		}
	}

	return mlswitch.Routes.Default
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

// Helper function to get the MLSwitch by its name
func getMLSwitchByName(name string, fullNetwork *[]importer.NetworkYAML) importer.MLSwitchYAML {
	for _, network := range *fullNetwork {
		for _, mls := range network.MLSwitches {
			if mls.Name == name {
				return mls
			}
		}
	}
	return importer.MLSwitchYAML{}
}

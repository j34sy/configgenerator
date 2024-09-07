package devices

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/datahandling"
	"github.com/j34sy/configgenerator/pkg/importer"
)

// Function to get the direct neighbors of each routing device
func GetDirectNeighbors(routing RoutingDevice, fullNetwork []RoutingDevice) map[string]string {
	neighbors := make(map[string]string)

	for _, device := range fullNetwork {
		if device.Name == routing.Name {
			continue
		}
		for _, iface := range device.Interfaces {
			for _, routerIface := range routing.Interfaces {
				check, err := datahandling.IsSameNetwork(routerIface.IP, iface.IP)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if check {
					neighbors[device.Name] = iface.IP
					fmt.Println("Found neighbor", device.Name, "for", routing.Name)
				}
			}
		}
	}

	return neighbors
}

// Function to find the next hop recursively
func FindNextHop(dest string, routing RoutingDevice, fullNetwork *[]importer.NetworkYAML) string {
	visited := make(map[string]bool)
	routingDevices := []RoutingDevice{}
	for _, network := range *fullNetwork {
		for _, router := range network.Routers {
			ifaces := []Interface{}
			for _, iface := range router.Interfaces {
				ifaces = append(ifaces, Interface{iface.Name, iface.Vlan, iface.IP, iface.Trunk, iface.Access, nil, iface.Native})
			}
			routingDevices = append(routingDevices, RoutingDevice{router.Name, ifaces, router.Routes.Destinations, router.Routes.Default})
		}
		for _, mlswitch := range network.MLSwitches {
			ifaces := []Interface{}
			for _, iface := range mlswitch.Interfaces {
				ifaces = append(ifaces, Interface{iface.Name, iface.Vlan, iface.IP, iface.Trunk, iface.Access, nil, iface.Native})
			}
			routingDevices = append(routingDevices, RoutingDevice{mlswitch.Name, ifaces, mlswitch.Routes.Destinations, mlswitch.Routes.Default})
		}
	}
	for _, device := range routingDevices {
		visited[device.Name] = false
	}
	success := false
	var nextHop string

	for !success {
		nextHop, success, visited = findNextHopRecursive(dest, routing, routingDevices, visited, false)
		fmt.Println("Next hop for ", dest, " is ", nextHop, " success: ", success)
		if success {
			return nextHop
		}
		success = true
		for _, device := range routingDevices {
			if !visited[device.Name] {
				success = false
				continue
			} else {
				continue
			}

		}

	}

	return nextHop
}

// Recursive function to find the next hop for routing devices
func findNextHopRecursive(dest string, routing RoutingDevice, fullNetwork []RoutingDevice, visited map[string]bool, remember bool) (string, bool, map[string]bool) {
	// Mark the current routing device as visited
	visited[routing.Name] = true

	neighbors := GetDirectNeighbors(routing, fullNetwork)

	// Check direct neighbors first
	for neighborName, neighborIP := range neighbors {
		neighborDevice := getRoutingDeviceByName(neighborName, fullNetwork)
		for _, iface := range neighborDevice.Interfaces {
			check, err := datahandling.IsSameNetwork(iface.IP, dest)
			if err != nil {
				fmt.Println(err)
				return routing.Default, remember, visited
			}
			if check {
				fmt.Println("Found direct route to", dest, "via", neighborName)
				return neighborIP, true, visited
			}
		}
	}

	// If no direct route is found, check for multi-hop routes
	for neighborName, neighborIP := range neighbors {
		if !visited[neighborName] {
			neighborDevice := getRoutingDeviceByName(neighborName, fullNetwork)

			nextHop, remember, _ := findNextHopRecursive(dest, neighborDevice, fullNetwork, visited, remember)

			if nextHop != routing.Default {
				fmt.Println("Found next hop", nextHop, "for", dest, "via", neighborName)
				return neighborIP, remember, visited
			}
		}
	}
	fmt.Println("No route found for", dest, " visited:", visited)
	return routing.Default, false, visited
}

// Helper function to get the routing device by its name
func getRoutingDeviceByName(name string, devices []RoutingDevice) RoutingDevice {
	for _, device := range devices {
		if device.Name == name {
			return device
		}
	}
	fmt.Println("Device not found")
	return RoutingDevice{}
}

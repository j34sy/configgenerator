package devices

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/datahandling"
	"github.com/j34sy/configgenerator/pkg/importer"
)

func FindNextHop(dest string, routingDev RoutingDevice, fullNetwork *[]importer.NetworkYAML) string {
	visited := make(map[string]bool)

	for _, iface := range routingDev.Interfaces {

		if iface.IP == "" {
			continue
		}
		//fmt.Println("Requesting IP adresses for dev ", routingDev.Name, " with IPs ", dest, iface.IP)
		match, err := datahandling.IsSameNetwork(dest, iface.IP)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if match {
			return iface.IP
		}
	}

	devices := []RoutingDevice{}

	for _, network := range *fullNetwork {
		for _, router := range network.Routers {
			devices = append(devices, RoutingDevice{Name: router.Name, Interfaces: convertYAMLInterfaceToInterface(router.Interfaces), Destinations: router.Routes.Destinations, Default: router.Routes.Default})
		}
		for _, mls := range network.MLSwitches {
			if mls.Routing {
				devices = append(devices, RoutingDevice{Name: mls.Name, Interfaces: convertYAMLInterfaceToInterface(mls.Interfaces), Destinations: mls.Routes.Destinations, Default: mls.Routes.Default})
			}
		}
	}

	for _, device := range devices {
		visited[device.Name] = false
		if device.Name == routingDev.Name {
			visited[device.Name] = true
		}
	}

	for i := 0; i < len(GetDirectNeighbors(routingDev, devices)); i++ {
		nextHop, found, _ := findNextHopRecursive(dest, routingDev, devices, visited)

		if found {
			//fmt.Println("Found next hop: ", nextHop, " for ", dest, " should be ", ip, " on ", neighbor)
			return nextHop
		} // else {}
		//fmt.Println("Could not find next hop")
		//fmt.Println("Tried to find next hop for ", dest, " on ", routingDev.Name, " with visited: ", visited)
		// }
	}
	//fmt.Println("Could not find next hop")
	return routingDev.Default
}

func findNextHopRecursive(dest string, routingDev RoutingDevice, devices []RoutingDevice, visited map[string]bool) (string, bool, map[string]bool) {
	//fmt.Println("Visiting device: ", routingDev.Name)

	for _, iface := range routingDev.Interfaces {
		if iface.IP == "" {
			continue
		}
		//fmt.Println("Requesting IP adresses for dev ", routingDev.Name, " with IPs ", dest, iface.IP)
		match, err := datahandling.IsSameNetwork(dest, iface.IP)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if match {
			//fmt.Println("Match found on interface: ", iface.IP)
			return iface.IP, true, visited
		}
	}

	neighbors := GetDirectNeighbors(routingDev, devices)

	for neighbor, ip := range neighbors {
		if visited[neighbor] {
			//fmt.Println("Already visited: ", neighbor)
			continue
		}
		visited[neighbor] = true

		neighborDev := getRoutingDeviceByName(neighbor, devices)

		_, found, visited := findNextHopRecursive(dest, neighborDev, devices, visited)

		if found {
			//fmt.Println("Next hop found via neighbor: ", neighbor, " with IP: ", ip)
			return ip, true, visited
		}
	}

	return "", false, visited
}

func GetDirectNeighbors(routing RoutingDevice, fullNetwork []RoutingDevice) map[string]string {
	neighbors := make(map[string]string)

	for _, device := range fullNetwork {
		if device.Name == routing.Name {
			continue
		}
		for _, iface := range device.Interfaces {
			for _, routerIface := range routing.Interfaces {
				if routerIface.IP != "" && iface.IP != "" {
					//fmt.Println("Requesting IP adresses for neighbor checking ", device.Name, " with IPs ", routerIface, iface.IP)
					check, err := datahandling.IsSameNetwork(routerIface.IP, iface.IP)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if check {
						neighbors[device.Name] = iface.IP
					}
				}
			}
		}
	}

	return neighbors
}

func getRoutingDeviceByName(name string, devices []RoutingDevice) RoutingDevice {
	for _, device := range devices {
		if device.Name == name {
			return device
		}
	}
	return RoutingDevice{}
}

func convertYAMLInterfaceToInterface(yamlInterface []importer.InterfaceYAML) []Interface {
	interfaces := []Interface{}
	for _, iface := range yamlInterface {
		interfaces = append(interfaces, Interface{Name: iface.Name, IP: iface.IP})
	}
	return interfaces
}

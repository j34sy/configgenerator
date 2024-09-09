package devices

import (
	"fmt"

	"github.com/j34sy/configgenerator/pkg/datahandling"
	"github.com/j34sy/configgenerator/pkg/importer"
)

func FindNextHopv6(dest string, routingDev RoutingDevice, fullNetwork *[]importer.NetworkYAML) string {
	visited := make(map[string]bool)

	for _, iface := range routingDev.Interfaces {

		if iface.IPv6 == "" {
			continue
		}

		match, err := datahandling.IsSameNetworkv6(dest, iface.IPv6)
		if err != nil {
			continue
		}
		if match {
			return iface.IPv6
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

	for i := 0; i < len(GetDirectNeighborsv6(routingDev, devices)); i++ {
		nextHop, found, _ := findNextHopRecursivev6(dest, routingDev, devices, visited)

		if found {
			return nextHop
		}
	}

	return routingDev.Defaultv6
}

func findNextHopRecursivev6(dest string, routingDev RoutingDevice, devices []RoutingDevice, visited map[string]bool) (string, bool, map[string]bool) {

	for _, iface := range routingDev.Interfaces {

		if iface.IPv6 == "" {
			continue
		}

		match, err := datahandling.IsSameNetworkv6(dest, iface.IPv6)
		if err != nil {
			continue
		}
		if match {
			return iface.IPv6, true, visited
		}
	}

	neighbors := GetDirectNeighborsv6(routingDev, devices)

	for neighbor, ipv6 := range neighbors {
		if visited[neighbor] {
			continue
		}
		visited[neighbor] = true

		neighborDev := getRoutingDeviceByName(neighbor, devices)

		_, found, visited := findNextHopRecursivev6(dest, neighborDev, devices, visited)

		if found {
			return ipv6, true, visited
		}
	}

	return "", false, visited
}

func GetDirectNeighborsv6(routingDev RoutingDevice, devices []RoutingDevice) map[string]string {
	neighbors := make(map[string]string)

	for _, device := range devices {
		if device.Name == routingDev.Name {
			continue
		}
		for _, iface := range device.Interfaces {
			for _, routerIface := range routingDev.Interfaces {
				if routerIface.IPv6 != "" && iface.IPv6 != "" {
					check, err := datahandling.IsSameNetworkv6(routerIface.IPv6, iface.IPv6)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if check {
						neighbors[device.Name] = iface.IPv6
					}
				}
			}
		}
	}

	return neighbors
}

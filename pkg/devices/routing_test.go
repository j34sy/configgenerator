package devices

import (
	"testing"

	"github.com/j34sy/configgenerator/pkg/importer"
)

func TestFindNextHop(t *testing.T) {
	// Mock data

	fullNetwork := []importer.NetworkYAML{
		{
			Routers: []importer.RouterYAML{
				{
					Name: "Router1",
					Routes: importer.RoutesYAML{
						Default:      "192.168.1.2",
						Destinations: []string{"192.168.2.0/24", "192.168.3.0/24", "192.168.4.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.1.1/24"},
					},
				},
				{
					Name: "Router2",
					Routes: importer.RoutesYAML{
						Default:      "192.168.2.2",
						Destinations: []string{"192.168.3.0/24", "192.168.4.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.2.1/24"},
						{Name: "eth1", IP: "192.168.1.2/24"},
					},
				},
				{
					Name: "Router3",
					Routes: importer.RoutesYAML{
						Default:      "192.168.3.2",
						Destinations: []string{"192.168.1.0/24", "192.168.4.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.3.1/24"},
						{Name: "eth1", IP: "192.168.2.2/24"},
					},
				},
				{
					Name: "Router4",
					Routes: importer.RoutesYAML{
						Default:      "192.168.4.2",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.4.1/24"},
						{Name: "eth1", IP: "10.0.0.0/24"},
					},
				},
			},
			MLSwitches: []importer.MLSwitchYAML{
				{
					Name: "MLSwitch1",
					Routes: importer.RoutesYAML{
						Default:      "192.168.5.1",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.4.2/24"},
						{Name: "eth1", IP: "192.168.5.1/24"},
					},
				},
				{
					Name: "MLSwitch2",
					Routes: importer.RoutesYAML{
						Default:      "",
						Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "192.168.4.0/24"},
					},
					Interfaces: []importer.InterfaceYAML{
						{Name: "eth0", IP: "192.168.5.2/24"},
						{Name: "eth1", IP: "172.16.0.0/24"},
					},
				},
			},
		},
	}

	tests := []struct {
		dest       string
		expected   string
		routingDev RoutingDevice
	}{
		{"192.168.1.0/24", "192.168.2.1/24", RoutingDevice{
			Name: "Router3",
			Interfaces: []Interface{
				{Name: "eth0", IP: "192.168.3.1/24"},
				{Name: "eth1", IP: "192.168.2.2/24"},
			},
			Destinations: []string{"192.168.1.0/24", "192.168.4.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24"},
			Default:      "192.168.3.2"}},
		//FIXME: somehow not correct
		{"192.168.5.2/24", "192.168.3.2/24", RoutingDevice{
			Name: "Router3",
			Interfaces: []Interface{
				{Name: "eth0", IP: "192.168.3.1/24"},
				{Name: "eth1", IP: "192.168.2.2/24"},
			},
			Destinations: []string{"192.168.1.0/24", "192.168.4.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24"},
			Default:      "192.168.3.2"}},
		//FIXME: somehow not correct
		{"172.16.0.0/24", "192.168.3.2/24", RoutingDevice{
			Name: "Router3",
			Interfaces: []Interface{
				{Name: "eth0", IP: "192.168.3.1/24"},
				{Name: "eth1", IP: "192.168.2.2/24"},
			},
			Destinations: []string{"192.168.1.0/24", "192.168.4.0/24", "192.168.5.0/24", "10.0.0.0/8", "172.16.0.0/24"},
			Default:      "192.168.3.2"}},
		{"10.0.0.0/8", "192.168.5.1/24", RoutingDevice{
			Name: "MLSwitch2",
			Interfaces: []Interface{
				{Name: "eth0", IP: "192.168.5.2/24"},
				{Name: "eth1", IP: "172.16.0.0/24"},
			},
			Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "192.168.4.0/24", "10.0.0.0/8"},
			Default:      ""}},
		{"172.16.0.0/24", "192.168.5.2/24", RoutingDevice{
			Name: "MLSwitch1",
			Interfaces: []Interface{
				{Name: "eth0", IP: "192.168.5.1/24"},
				{Name: "eth1", IP: "192.168.4.2/24"},
			},
			Destinations: []string{"192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24", "10.0.0.0/8", "172.16.0.0/24"},
			Default:      "192.168.5.2"}},
	}

	for _, test := range tests {
		t.Run(test.dest, func(t *testing.T) {
			result := FindNextHop(test.dest, test.routingDev, &fullNetwork)
			if result != test.expected {
				t.Errorf("FindNextHop(%s) = %s; expected %s", test.dest, result, test.expected)
			}
		})
	}
}

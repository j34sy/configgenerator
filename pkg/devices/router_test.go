package devices

import (
	"testing"

	"github.com/j34sy/configgenerator/pkg/importer"
)

func TestFindNextHop(t *testing.T) {
	// Mock data
	router := importer.RouterYAML{
		Name: "Router3",
		Routes: importer.RoutesYAML{
			Default:      "192.168.3.2",
			Destinations: []string{"192.168.1.0/24", "192.168.4.0/24"},
		},
		Interfaces: []importer.InterfaceYAML{
			{Name: "eth0", IP: "192.168.3.1/24"},
			{Name: "eth1", IP: "192.168.2.2/24"},
		},
	}

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
			},
		},
	}

	tests := []struct {
		dest     string
		expected string
		router   importer.RouterYAML
	}{
		{"192.168.1.0/24", "192.168.3.1/24", router},
		{"0.0.0.0/24", "192.168.4"},
		{"172.16.1.5", "172.16.0.1"},
		{"192.168.2.5", "192.168.2.1"},
		{"8.8.8.8", "0.0.0.0"},
	}

	for _, test := range tests {
		t.Run(test.dest, func(t *testing.T) {
			result := FindNextHop(test.dest, router, &fullNetwork)
			if result != test.expected {
				t.Errorf("FindNextHop(%s) = %s; expected %s", test.dest, result, test.expected)
			}
		})
	}
}

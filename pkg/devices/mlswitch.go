package devices

import (
	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateMLSwitch(mlSwitchYAML importer.MLSwitchYAML) *MLSwitch {
	return &MLSwitch{
		RoutingDevice: RoutingDevice{
			mlSwitchYAML.Name,
			[]Interface{},
			mlSwitchYAML.Routes.Destinations,
			mlSwitchYAML.Routes.Default,
		},
		OSPFRouters: []OSPFRouter{},
		Routing:     mlSwitchYAML.Routing,
		Users:       []User{},
		Routes:      []Route{},
	}
}

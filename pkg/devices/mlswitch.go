package devices

import (
	"github.com/j34sy/configgenerator/pkg/importer"
)

func CreateMLSwitch(mlSwitchYAML importer.MLSwitchYAML) *MLSwitch {
	return &MLSwitch{
		Name:        mlSwitchYAML.Name,
		Interfaces:  []Interface{},
		Routes:      []Route{},
		OSPFRouters: []OSPFRouter{},
		Routing:     mlSwitchYAML.Routing,
		Users:       []User{},
	}
}

package importer

type User struct {
	Name      string `yaml:"name"`
	Password  string `yaml:"password"`
	Privilege string `yaml:"privilege"`
}

type Vlan struct {
	ID      int    `yaml:"id"`
	Name    string `yaml:"name"`
	Subnet  string `yaml:"subnet,omitempty"`
	Gateway string `yaml:"gateway,omitempty"`
}

type VlanGroup struct {
	Switches []string `yaml:"switches"`
	List     []Vlan   `yaml:"list"`
}

type Interface struct {
	Name   string `yaml:"name"`
	Vlan   string `yaml:"vlan,omitempty"`
	IP     string `yaml:"ip,omitempty"`
	Trunk  []int  `yaml:"trunk,omitempty"`
	Access int    `yaml:"access,omitempty"`
	OSPF   *OSPF  `yaml:"ospf,omitempty"`
}

type OSPF struct {
	Process int `yaml:"process"`
	Area    int `yaml:"area"`
}

type Routes struct {
	Destinations []string `yaml:"destinations,omitempty"`
	Default      string   `yaml:"default,omitempty"`
}

type OSPFRouter struct {
	Process int    `yaml:"process"`
	ID      string `yaml:"id"`
}

type Router struct {
	Name       string       `yaml:"name"`
	Interfaces []Interface  `yaml:"interfaces,omitempty"`
	Routes     Routes       `yaml:"routes,omitempty"`
	OSPFRouter []OSPFRouter `yaml:"ospf,omitempty"`
}

type Switch struct {
	Name       string      `yaml:"name"`
	Interfaces []Interface `yaml:"interfaces,omitempty"`
}

type MLSwitch struct {
	Name       string       `yaml:"name"`
	Routing    bool         `yaml:"routing,omitempty"`
	Interfaces []Interface  `yaml:"interfaces,omitempty"`
	Routes     Routes       `yaml:"routes,omitempty"`
	OSPFRouter []OSPFRouter `yaml:"ospf,omitempty"`
}

type Network struct {
	Name       string      `yaml:"name"`
	Users      []User      `yaml:"users,omitempty"`
	Vlans      []VlanGroup `yaml:"vlans,omitempty"`
	Routers    []Router    `yaml:"routers,omitempty"`
	Switches   []Switch    `yaml:"switches,omitempty"`
	MLSwitches []MLSwitch  `yaml:"mlswitches,omitempty"`
}

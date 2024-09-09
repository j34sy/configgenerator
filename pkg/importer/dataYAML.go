package importer

type UserYAML struct {
	Name      string `yaml:"name"`
	Password  string `yaml:"password"`
	Privilege string `yaml:"privilege"`
}

type VlanYAML struct {
	ID        int    `yaml:"id"`
	Name      string `yaml:"name"`
	Subnet    string `yaml:"subnet,omitempty"`
	Gateway   string `yaml:"gateway,omitempty"`
	Subnetv6  string `yaml:"subnetv6,omitempty"`
	Gatewayv6 string `yaml:"gatewayv6,omitempty"`
}

type VlanGroupYAML struct {
	Switches []string   `yaml:"switches"`
	List     []VlanYAML `yaml:"list"`
}

type InterfaceYAML struct {
	Name   string    `yaml:"name"`
	Vlan   string    `yaml:"vlan,omitempty"`
	IP     string    `yaml:"ip,omitempty"`
	IPv6   string    `yaml:"ipv6,omitempty"`
	Trunk  []int     `yaml:"trunk,omitempty"`
	Access int       `yaml:"access,omitempty"`
	OSPF   *OSPFYAML `yaml:"ospf,omitempty"`
	Native int       `yaml:"native,omitempty"`
}

type OSPFYAML struct {
	Process int `yaml:"process"`
	Area    int `yaml:"area"`
}

type RoutesYAML struct {
	Destinations   []string `yaml:"destinations,omitempty"`
	Default        string   `yaml:"default,omitempty"`
	Destinationsv6 []string `yaml:"destinationsv6,omitempty"`
	Defaultv6      string   `yaml:"defaultv6,omitempty"`
}

type OSPFRouterYAML struct {
	Process int    `yaml:"process"`
	ID      string `yaml:"id"`
}

type RouterYAML struct {
	Name       string           `yaml:"name"`
	Interfaces []InterfaceYAML  `yaml:"interfaces,omitempty"`
	Routes     RoutesYAML       `yaml:"routes,omitempty"`
	OSPFRouter []OSPFRouterYAML `yaml:"ospf,omitempty"`
}

type SwitchYAML struct {
	Name       string          `yaml:"name"`
	Interfaces []InterfaceYAML `yaml:"interfaces,omitempty"`
}

type MLSwitchYAML struct {
	Name       string           `yaml:"name"`
	Routing    bool             `yaml:"routing,omitempty"`
	Interfaces []InterfaceYAML  `yaml:"interfaces,omitempty"`
	Routes     RoutesYAML       `yaml:"routes,omitempty"`
	OSPFRouter []OSPFRouterYAML `yaml:"ospf,omitempty"`
}

type NetworkYAML struct {
	Name         string          `yaml:"name"`
	EnableSecret string          `yaml:"enable,omitempty"`
	Users        []UserYAML      `yaml:"users,omitempty"`
	Vlans        []VlanGroupYAML `yaml:"vlans,omitempty"`
	Routers      []RouterYAML    `yaml:"routers,omitempty"`
	Switches     []SwitchYAML    `yaml:"switches,omitempty"`
	MLSwitches   []MLSwitchYAML  `yaml:"mlswitches,omitempty"`
}

package devices

type RoutingDevice struct {
	Name           string
	Interfaces     []Interface
	Destinations   []string
	Destinationsv6 []string
	Default        string
	Defaultv6      string
}

type Router struct {
	RoutingDevice
	Domain       string
	EnableSecret string
	OSPFRouters  []OSPFRouter
	Users        []User
	Routes       []Route
	Routesv6     []Routev6
}

type Switch struct {
	Name         string
	Domain       string
	EnableSecret string
	Interfaces   []Interface
	Vlans        []Vlan
	Users        []User
	Default      string
	Defaultv6    string
}

type MLSwitch struct {
	RoutingDevice
	Domain       string
	EnableSecret string
	Vlans        []Vlan
	OSPFRouters  []OSPFRouter
	Routing      bool
	Users        []User
	Routes       []Route
	Routesv6     []Routev6
}

type User struct {
	Name      string
	Password  string
	Privilege string
}

type Vlan struct {
	ID        int
	Name      string
	Subnet    string
	Gateway   string
	Subnetv6  string
	Gatewayv6 string
}

type Interface struct {
	Name   string
	Vlan   string
	IP     string
	IPv6   string
	Trunk  []int
	Access int
	OSPF   *OSPF
	Native int
}

type OSPF struct {
	Process int
	Area    int
}

type OSPFRouter struct {
	Process int
	ID      string
}

type Route struct {
	Destination string
	NextHop     string
}

type Routev6 struct {
	Destination string
	NextHop     string
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

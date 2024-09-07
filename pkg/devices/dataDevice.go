package devices

type Router struct {
	Name        string
	Domain      string
	Interfaces  []Interface
	Routes      []Route
	OSPFRouters []OSPFRouter
	Users       []User
	Default     string
}

type Switch struct {
	Name       string
	Domain     string
	Interfaces []Interface
	Vlans      []Vlan
	Users      []User
	Default    string
}

type MLSwitch struct {
	Name        string
	Domain      string
	Interfaces  []Interface
	Vlans       []Vlan
	Routes      []Route
	OSPFRouters []OSPFRouter
	Routing     bool
	Users       []User
	Default     string
}

type User struct {
	Name      string
	Password  string
	Privilege string
}

type Vlan struct {
	ID      int
	Name    string
	Subnet  string
	Gateway string
}

type Interface struct {
	Name   string
	Vlan   string
	IP     string
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

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

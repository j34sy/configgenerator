package datahandling

import (
	"strconv"
	"strings"

	"github.com/j34sy/SubnetCalculator/pkg/subnetcalc"
)

type ErrInvalidIPv4Address struct {
	Message string
}

func (e *ErrInvalidIPv4Address) Error() string {
	return e.Message
}

// Fetch IPv4 address data
func GetIPv4Address(ipv4String string) (*subnetcalc.IPv4Address, error) {
	var ip [4]int
	var cidr uint8

	slashSplit := strings.Split(ipv4String, "/")
	if len(slashSplit) != 2 {
		return nil, &ErrInvalidIPv4Address{Message: "Invalid IPv4 address format, needs to be 1.2.3.4/12"}
	}

	ipSplit := strings.Split(slashSplit[0], ".")
	if len(ipSplit) != 4 {
		return nil, &ErrInvalidIPv4Address{Message: "Invalid IPv4 address format, needs to be 1.2.3.4/12"}
	}

	cidrInt, err := strconv.Atoi(slashSplit[1])

	if err != nil {
		return nil, &ErrInvalidIPv4Address{Message: "Invalid CIDR format, needs to be 1.2.3.4/12"}
	}

	if cidrInt > 0 && cidrInt < 33 {
		cidr = uint8(cidrInt)
	} else {
		return nil, &ErrInvalidIPv4Address{Message: "Invalid CIDR format, needs to be 1.2.3.4/12; CIDR must be between 0 and 32"}
	}

	for i, octet := range ipSplit {
		octetInt, err := strconv.Atoi(octet)
		if err != nil {
			return nil, &ErrInvalidIPv4Address{Message: "Invalid IPv4 address format, needs to be 1.2.3.4/12"}
		}
		if octetInt >= 0 && octetInt <= 255 {
			ip[i] = octetInt
		} else {
			return nil, &ErrInvalidIPv4Address{Message: "Invalid IPv4 address format, needs to be 1.2.3.4/12; octet must be between 0 and 255"}
		}
	}
	rawIP := subnetcalc.NewIPv4Address(ip, cidr)
	rawIP.Calculate()
	return rawIP, nil
}

func IsSameNetwork(ipv4AString string, ipv4BString string) (bool, error) {
	ipv4A, err := GetIPv4Address(ipv4AString)
	if err != nil {
		return false, err
	}

	ipv4B, err := GetIPv4Address(ipv4BString)
	if err != nil {
		return false, err
	}

	result := false

	if ipv4A.GetNetworkAddress() == ipv4B.GetNetworkAddress() {
		if ipv4A.GetBroadcastAddress() == ipv4B.GetBroadcastAddress() {
			result = true
		}
	}

	return result, nil
}

package datahandling

import (
	"net"
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
		return nil, &ErrInvalidIPv4Address{Message: "Invalid IPv4 address format, needs to be 1.2.3.4/12, got: " + ipv4String}
	}

	ipSplit := strings.Split(slashSplit[0], ".")
	if len(ipSplit) != 4 {
		return nil, &ErrInvalidIPv4Address{Message: "Invalid IPv4 address format, needs to be 1.2.3.4/12, got: " + ipv4String}
	}

	cidrInt, err := strconv.Atoi(slashSplit[1])

	if err != nil {
		return nil, &ErrInvalidIPv4Address{Message: "Invalid CIDR format, needs to be 1.2.3.4/12, got: " + ipv4String}
	}

	if cidrInt > 0 && cidrInt < 33 {
		cidr = uint8(cidrInt)
	} else {
		return nil, &ErrInvalidIPv4Address{Message: "Invalid CIDR format, needs to be 1.2.3.4/12; CIDR must be between 0 and 32"}
	}

	for i, octet := range ipSplit {
		octetInt, err := strconv.Atoi(octet)
		if err != nil {
			return nil, &ErrInvalidIPv4Address{Message: "Invalid IPv4 address format, needs to be 1.2.3.4/12, got: " + ipv4String}
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

func IsSameNetworkv6(ipv6A string, ipv6B string) (bool, error) {
	v6A, _, err1 := net.ParseCIDR(ipv6A)
	v6B, _, err2 := net.ParseCIDR(ipv6B)

	if err1 != nil {
		return false, err1
	}
	if err2 != nil {
		return false, err2
	}

	subnetA := v6A.Mask(net.CIDRMask(64, 128))
	subnetB := v6B.Mask(net.CIDRMask(64, 128))

	return subnetA.Equal(subnetB), nil
}

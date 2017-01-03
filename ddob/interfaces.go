package ddob

import (
	"net"
	"regexp"
)

func ListInterfaceAddrs(namePattern *regexp.Regexp) (addrs []net.Addr, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, inter := range interfaces {
		if !namePattern.MatchString(inter.Name) {
			continue
		}

		interfaceAddrs, err := inter.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range interfaceAddrs {
			addrs = append(addrs, addr)
		}
	}

	return
}

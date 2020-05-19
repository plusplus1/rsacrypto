package utilLib

import (
	"net"
)

import (
	"parallel_rsa/commonLib/logLib"
)

var (
	log = logLib.Std
)

func GetLocalhostIPV4() string {
	arrAddress, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("GetLocalhostIPV4 fail, %v", err)
		return ""
	}
	for _, address := range arrAddress {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

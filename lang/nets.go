package lang

import (
	"fmt"
	"net"
)

func MacAddrToInt(macByte []byte, maxDataCenterId int64) int64 {
	macByteLen := len(macByte)
	if macByteLen == 0 {
		return 1
	}
	id := int64(0)
	id = ((0x000000FF & int64(macByte[macByteLen-1])) | (0x0000FF00 & int64(macByte[macByteLen-2]) << 8)) >> 6
	id = id % maxDataCenterId
	return id
}

func GetMacAddr() []byte {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return []byte{}
	}

	for _, netInterface := range netInterfaces {
		if netInterface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if netInterface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := netInterface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			mac := netInterface.HardwareAddr
			if len(mac) == 0 {
				continue
			}
			return mac
		}
	}
	return []byte{}
}

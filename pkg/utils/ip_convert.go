package utils

import (
	"fmt"
	"math/big"
	"net"
)

func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func IP6toInt(ipv6 net.IP) *big.Int {
	ipv6int := big.NewInt(0)
	ipv6int.SetBytes(ipv6.To16())
	return ipv6int
}

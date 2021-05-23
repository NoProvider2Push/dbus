package main

import (
	"net"

	"github.com/karmanyaahm/np2p_linux/storage"
)

func getYggIps() (ips []net.IP) {
	_, rang, _ := net.ParseCIDR("200::/8")
	inter, _ := net.Interfaces()

	for _, i := range inter {
		addrs, _ := i.Addrs()
		for _, j := range addrs {
			ip, _, _ := net.ParseCIDR(j.String())
			if rang.Contains(ip) {
				ips = append(ips, ip)
			}
		}
	}
	return
}

func main() {
	//30043
	storage.InitDB()

}

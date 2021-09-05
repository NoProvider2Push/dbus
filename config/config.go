package config

import (
	"log"
	"net"
	"strconv"

	"github.com/karmanyaahm/np2p_linux/utils"
	"gopkg.in/ini.v1"
)

func GetEndpointURL() string {
	return *c.ProxyURL + "/" + GetIPPort()
}

func GetIPPort() string {
	return net.JoinHostPort(*c.IP, strconv.Itoa(*c.Port))
}

var c conf

type conf struct {
	IP   *string
	Port *int

	ProxyURL *string
}

//                  public id - private id
type application = map[string]string

func Init(name string) {
	cfg, err := ini.Load(utils.StoragePath(name + ".conf"))
	if err != nil {
		log.Fatal(err)
	}

	c = conf{}
	cfg.Section("").MapTo(c)

	secs := cfg.Sections()
	secs[0].MapTo(&c)

	defaults(&c)

	_ = secs[0].ReflectFrom(&c)
	_ = cfg.SaveTo(utils.StoragePath(name + ".conf"))
}

func defaults(c *conf) {
	if c.IP == nil || ipExists(*c.IP) {
		ygg := getYggIps()
		if len(ygg) > 0 {
			ip := ygg[0].To16().String()
			c.IP = &ip
		} else {
			log.Fatalln("No IP on this machine available in config")
		}
	}

	if c.Port == nil {
		p := 30043
		c.Port = &p
	}

	if c.ProxyURL == nil {
		log.Fatalln("Need a proxy url, TODO setup a public one")
	}
}
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

func ipExists(s string) bool {
	inter, _ := net.Interfaces()

	for _, i := range inter {
		addrs, _ := i.Addrs()
		for _, j := range addrs {
			_, c, _ := net.ParseCIDR(j.String())
			if c.Contains(net.ParseIP(s)) {
				return true
			}
		}
	}
	return false
}

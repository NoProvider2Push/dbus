# No Provider 2 Push D-Bus

Config:

```ini
proxyurl = direct
# for testing or if your computer is publicly exposed to a static IP for some reason (still not recommended because np2p doesn't support tls(https))
# or
proxyurl = https://mynp2p.proxy.tld

port = 30043
# defaults to this so no need to fill in unless you want to change it

IP = 192.168.0.99
# ipv4
IP = 2001:0DB8::123
#ipv6

# depends on your proxy setup
# this defaults to a Yggdrasil IP address if you're running that in the background
IP = 201:be::0123
```


Roadmap: 
- alpha: currently
- beta: once builds are set up
- stable: v1.0 should be released once dbus UP platform is proven stable - don't know timeline

## Library

The distributor package can be used as a module in your own distributor if you wish. Other parts like config and storage can also be copied with the appropriate license. NP2P is usually the 'example distributor' in UnifiedPush due to its simplicity.

```sh
go get -u unifiedpush.org/go/np2p/dbus # for the np2p parts
go get -u unifiedpush.org/go/dbus/distributor # for the distributor API (this one is more applicable to other distributors than the one above)
```

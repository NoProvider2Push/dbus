package distributor

import "github.com/godbus/dbus"

type DBus struct {
	client *dbus.Conn
}

func (d DBus) NewDistributor(distName string) *Distributor {
	//register on bus
	return &Distributor{
		name: distName,
	}
}

type Distributor struct {
	name string
}

func (d Distributor) Register(name, token string) (thing, reason string, err *dbus.Error) {

	return "REGISTRATION_FAILED", "unimpl", nil
}

func (d Distributor) Unregister(token string) *dbus.Error {
	return nil
}

func (d DBus) NewConector(objstr string) *Connector {
	obj := d.client.Object(objstr, "/org/unifiedpush/Connector")
	return &Connector{
		obj: &obj,
	}
}

type Connector struct {
	obj *dbus.BusObject
}

func (c Connector) Message(token, contents, id string) error {
	return nil
}

func (c Connector) NewEndpoint(token, endpoint string) error {
	return nil
}

func (c Connector) Unregistered(token string) error {
	return nil
}

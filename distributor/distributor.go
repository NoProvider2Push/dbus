package distributor

import "github.com/godbus/dbus"

type DBus struct {
	client *dbus.Conn
}

func (d DBus) NewDistributor(distName string) *Distributor {
	//register on bus
	return &Distributor{
		client: d,
		name:   distName,
	}
}

type Distributor struct {
	client DBus
	name   string
}

func (d Distributor) Register(appid, token string) (thing, reason string, err *dbus.Error) {
	c := d.client.NewConnector(appid)
	if err := c.NewEndpoint(token, "http://.../UP?token="+token); err != nil {
		return "REGISTRATION_FAILED", err.Error(), nil
	}
	return "NEW_ENDPOINT", "", nil
}

func (d Distributor) Unregister(token string) *dbus.Error {
	return nil
}

func (d DBus) NewConector(appid string) *Connector {
	obj := d.client.Object(appid, "/org/unifiedpush/Connector")
	return &Connector{
		obj: &obj,
	}
}

type Connector struct {
	obj *dbus.BusObject
}

func (c Connector) Message(token, contents, id string) error {
	return c.obj.Call("org.unifiedpush.Connector1.Message", dbus.FlagNoReplyExpected, token, contents, id).Err
}

func (c Connector) NewEndpoint(token, endpoint string) error {
	return c.obj.Call("org.unifiedpush.Connector1.NewEndpoint", dbus.FlagNoReplyExpected, token, endpoint).Err
}

func (c Connector) Unregistered(token string) error {
	return c.obj.Call("org.unifiedpush.Connector1.Unregistered", dbus.FlagNoReplyExpected, token).Err
}

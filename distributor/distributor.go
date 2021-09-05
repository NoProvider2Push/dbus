package distributor

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

type Distrib interface {
	Register(appName, token string) (endpoint, refuseReason string, err error)
	Unregister(token string)
}

func NewDBus(distName string) *DBus {
	//register on bus
	return &DBus{
		name: distName,
	}
}

type DBus struct {
	client *dbus.Conn
	name   string
}

// StartHandling exports the distributor interface and requests the app's name on the bus
func (d *DBus) StartHandling(handler Distrib) (err error) {

	d.client, err = dbus.ConnectSessionBus()
	if err != nil {
		return err
	}

	err = d.client.Export(&dBusDistrib{handler: handler}, "/org/unifiedpush/Distributor", "org.unifiedpush.Distributor1")
	if err != nil {
		return err
	}

	name, err := d.client.RequestName(d.name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}
	if name != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("Cannot request name on dbus")
	}

	return nil
}

func (d DBus) Close() error {
	return nil
}

func (d DBus) NewConector(appid string) *Connector {
	obj := d.client.Object(appid, "/org/unifiedpush/Connector")
	return &Connector{
		obj: obj,
	}
}

type dBusDistrib struct {
	handler Distrib
}

func (d dBusDistrib) Register(appid, token string) (thing, reason string, err *dbus.Error) {
	endpoint, refused, errr := d.handler.Register(appid, token)
	if errr != nil {
		return "REGISTRATION_FAILED", errr.Error(), nil
	}
	if refused != "" {
		return "REGISTRATION_REFUSED", refused, nil
	}
	return "NEW_ENDPOINT", endpoint, nil
}

func (d dBusDistrib) Unregister(token string) *dbus.Error {
	d.handler.Unregister(token)
	return nil
}

type Connector struct {
	obj dbus.BusObject
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

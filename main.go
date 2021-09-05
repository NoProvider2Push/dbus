package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/karmanyaahm/np2p_linux/config"
	"github.com/karmanyaahm/np2p_linux/distributor"
	"github.com/karmanyaahm/np2p_linux/storage"
)

var store *storage.Storage
var dbus *distributor.DBus

func main() {
	store = storage.InitStorage("np2p")
	config.Init("np2p")

	dbus = distributor.NewDBus("org.unifiedpush.distributor.NP2P")

	dbus.StartHandling(handler{})

	http.HandleFunc("/", httpHandle)
	log.Fatal(http.ListenAndServe(config.GetIPPort(), nil))
}

func httpHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, `{"unifiedpush" : {"version" : 1}}`)
	} else if r.Method == http.MethodPost {
		parts := strings.Split(r.URL.Path, "/")
		var token string
		if len(parts) > 0 {
			token = parts[0]
		} else {
			w.WriteHeader(400)
			return
		}

		conn := store.GetConnectionbyPublic(token)
		if conn == nil {
			w.WriteHeader(404)
			return
		}

		body, _ := io.ReadAll(io.LimitReader(r.Body, 4005))
		if len(body) > 4003 {
			w.WriteHeader(413)
			return
		}

		w.WriteHeader(202)
		//implement 429 counter

		_ = dbus.NewConnector(conn.AppID).Message(conn.AppToken, string(body), "") //TODO errors

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type handler struct {
}

func (h handler) Register(appName, token string) (endpoint, refuseReason string, err error) {
	conn := store.NewConnection(appName, token)
	if conn != nil {
		return config.GetEndpointURL() + "?token=" + conn.PublicToken, "", nil
	}
	//np2p doesn't have a situation for refuse
	return "", "", errors.New("Unknown error with NoProvider2Push")
}
func (h handler) Unregister(token string) {
	deletedConn, err := store.DeleteConnection(token)
	if err != nil {
		//?????
	}
	_ = dbus.NewConnector(deletedConn.AppID).Unregistered(deletedConn.AppToken)
}

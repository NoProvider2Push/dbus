package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"unifiedpush.org/go/np2p_dbus/config"
	"unifiedpush.org/go/np2p_dbus/distributor"
	"unifiedpush.org/go/np2p_dbus/storage"
	"unifiedpush.org/go/np2p_dbus/utils"
)

var store *storage.Storage
var dbus *distributor.DBus

func main() {
	var err error
	if store, err = storage.InitStorage(utils.StoragePath("np2p.db")); err != nil {
		log.Fatalln("failed to connect database")
	}
	config.Init("np2p")

	dbus = distributor.NewDBus("org.unifiedpush.Distributor.NP2P")

	dbus.StartHandling(handler{})

	go handleEndpointSettingsChanges()

	http.HandleFunc("/", httpHandle)
	utils.Log.Debugln("listening on", config.GetIPPort(), "with endpoints like", config.GetEndpointURL("<token>"), "...")
	log.Fatal(http.ListenAndServe(config.GetIPPort(), nil))
}

func handleEndpointSettingsChanges() {
	endpointFormat := config.GetEndpointURL("<token>")
	for _, i := range store.GetUnequalSettings(endpointFormat) {
		utils.Log.Debugln("new endpoint format for", i.AppID, i.AppToken)
		//newconnection updates the endpoint settings when one already exists
		n := store.NewConnection(i.AppID, i.AppToken, endpointFormat)
		if n == nil || n.Settings != endpointFormat {
			utils.Log.Debugln("unable to save new endpoint format for", i.AppID, i.AppToken)
			continue
		}
		dbus.NewConnector(n.AppID).NewEndpoint(n.AppToken, config.GetEndpointURL(n.PublicToken))
	}
}

func httpHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, `{"unifiedpush" : {"version" : 1}}`)
	} else if r.Method == http.MethodPost {
		parts := strings.Split(r.URL.Path, "/")
		utils.Log.Debugln("received request from", r.URL.Path)

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

		_ = dbus.NewConnector(conn.AppID).Message(conn.AppToken, body, "") //TODO errors
		utils.Log.Infoln("MESSAGE", conn.AppID, conn.AppToken, "from", r.RemoteAddr)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type handler struct {
}

func (h handler) Register(appName, token string) (endpoint, refuseReason string, err error) {
	conn := store.NewConnection(appName, token, config.GetEndpointURL("<token>"))
	utils.Log.Debugln("registered", conn)
	if conn != nil {
		return config.GetEndpointURL(conn.PublicToken), "", nil
	}
	//np2p doesn't have a situation for refuse
	return "", "", errors.New("Unknown error with NoProvider2Push")
}
func (h handler) Unregister(token string) {
	deletedConn, err := store.DeleteConnection(token)
	utils.Log.Debugln("deleted", deletedConn)

	if err != nil {
		//?????
	}
	_ = dbus.NewConnector(deletedConn.AppID).Unregistered(deletedConn.AppToken)
}

// Package wsServer list keys for admin page
package wsServer

import (
	"github.com/Atluss/TestTaskElma/pkg/v1/api"
	"github.com/Atluss/TestTaskElma/pkg/v1/auth"
	"github.com/Atluss/TestTaskElma/pkg/v1/config"
	"github.com/Atluss/TestTaskElma/pkg/v1/dataKeys"
	"log"
	"net/http"
)

// WSList setup websocket connection for admin list client
func WSList(set *config.Setup, secure bool) error {
	wsClient := &wSList{
		Url:    "/ws_list",
		Setup:  set,
		Secure: secure, // only for admins!
	}
	set.Route.HandleFunc(wsClient.Url, wsClient.HandleConnection)
	return nil
}

type wSList struct {
	api.RequestRest
	Url    string
	Setup  *config.Setup
	Secure bool
}

type wSReplyList struct {
	Status int
	Type   string // message type
	// getList
	// updateKey
	// newKey
	Items      []dataKeys.Keys
	UpdatedKey string
}

type wSRequestList struct {
	Status int
	Type   string // message type
	// getList
	// updateKey
	Key  string
	Name string
}

// HandleConnection admins connections
func (obj *wSList) HandleConnection(w http.ResponseWriter, r *http.Request) {

	if obj.Secure { // check it's admin
		if !auth.CheckAuth(auth.GetSession(r)) {
			log.Printf("warning: request of rights from ip %s ", r.RemoteAddr)
			return
		}
	}

	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	admin := admin{
		Conn: ws,
		send: make(chan interface{}),
	}
	admin.Run(obj.Setup.Gorm)
}

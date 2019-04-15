// list keys for admin page
package ws_server

import (
	"github.com/Atluss/TestTaskElma/lib/api"
	"github.com/Atluss/TestTaskElma/lib/auth"
	"github.com/Atluss/TestTaskElma/lib/config"
	"github.com/Atluss/TestTaskElma/lib/data"
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
	Items      []data.Keys
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

	admin := Admin{
		Conn: ws,
		send: make(chan interface{}),
	}
	admin.Run(obj.Setup.Gorm)
}

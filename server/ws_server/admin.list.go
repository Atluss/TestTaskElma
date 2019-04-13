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

func WSList(set *config.Setup, secure bool) error {

	wsClient := &wSList{
		Url:    "/ws_list",
		Setup:  set,
		Secure: secure,
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
	Items  []data.Keys
}

type wSRequestList struct {
	Status int
}

func (obj *wSList) HandleConnection(w http.ResponseWriter, r *http.Request) {

	if obj.Secure {
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

	for {

		req := wSRequestList{}
		err = ws.ReadJSON(&req)
		if err != nil {
			log.Printf("error: %v", err)
			ws.CloseHandler()
			break
		}

		msg := wSReplyList{}
		if keys, err := data.GetKeysByStatus(req.Status, obj.Setup.Gorm); err != nil {
			msg.Status = http.StatusInternalServerError
		} else {
			msg.Status = http.StatusOK
			msg.Items = keys
		}

		err := ws.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			ws.CloseHandler()
			break
		}
	}

}

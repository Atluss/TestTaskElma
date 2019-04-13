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
	Type   string // message type
	// getList
	// updateKey
	Items      []data.Keys
	UpdatedKey string
}

type wSRequestList struct {
	Type string // message type
	// getList
	// updateKey
	Status int
	Key    string
	Name   string
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

		log.Printf("%+v", req)

		msg := wSReplyList{
			Status: http.StatusOK,
			Items:  []data.Keys{},
		}

		if req.Type == "getList" {
			msg.Type = "getList"
			msg.Items, err = data.GetKeysByStatus(req.Status, obj.Setup.Gorm)
			if err != nil {
				msg.Status = http.StatusInternalServerError
			}
		} else {

			msg.Type = "updateKey"

			key := data.Keys{
				Key:    req.Key,
				Name:   req.Name,
				Status: req.Status,
			}

			if err := key.Update(obj.Setup.Gorm); err != nil {
				msg.Status = http.StatusInternalServerError
			} else {
				msg.UpdatedKey = key.Key
			}

			log.Printf("%+v", key)
			log.Printf("%+v", msg)

		}

		err := ws.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			ws.CloseHandler()
			break
		}
	}

}

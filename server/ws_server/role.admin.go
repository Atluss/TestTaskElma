// websocket connection for admins
package ws_server

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib/data"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

var Admins = make(map[*Admin]bool) // admins connects

// handleNewKeyBroadcast send CPU load all accepted connections
func AddAllAdminMessage(msg interface{}) {
	for admin := range Admins {
		admin.send <- msg
	}
}

type Admin struct {
	Conn *websocket.Conn
	send chan interface{}
	run  chan string
}

// Run admin handlers
func (obj *Admin) Run(db *gorm.DB) {
	Admins[obj] = true
	go obj.Read(db)
	go obj.Write()
	wait := <-obj.run
	log.Println(wait)
}

// Stop connection and remove admin from admins list
func (obj *Admin) Stop(err error) {
	if obj.Conn != nil {
		obj.Conn.Close()
	}
	delete(Admins, obj)
	log.Printf("admins: len: %d %+v", len(Admins), Admins)
	obj.run <- fmt.Sprintf("stop: err %s", err)
}

// Write send messages
func (obj *Admin) Write() {
	for msg := range obj.send {
		err := obj.Conn.WriteJSON(msg)
		if err != nil {
			obj.Stop(err)
			break
		}
	}
}

// Read process input messages
func (obj *Admin) Read(db *gorm.DB) {
	for {
		req := wSRequestList{}
		err := obj.Conn.ReadJSON(&req)
		if err != nil {
			obj.Stop(err)
			break
		}

		log.Printf("%+v", req)
		msg := wSReplyList{
			Status: http.StatusOK,
			Items:  []data.Keys{},
		}

		switch req.Type {
		case GetList:
			{
				msg.Type = GetList
				if req.Status > 10 { // for online clients
					msg.Items = GetActiveClients()
				} else {
					msg.Items, err = data.GetKeysByStatus(req.Status, db)
					if err != nil {
						msg.Status = http.StatusInternalServerError
					}
				}
				obj.send <- msg
			}
		case UpdateKey:
			{
				msg.Type = UpdateKey
				key := data.Keys{
					Key:    req.Key,
					Name:   req.Name,
					Status: req.Status,
				}

				if err := key.Update(db); err != nil {
					msg.Status = http.StatusInternalServerError
				} else {
					msg.UpdatedKey = key.Key
					if key.Status == data.KeyBlocked {
						BanKeys.V[key.Key] = true
					} else {
						delete(BanKeys.V, key.Key)
					}
				}

				stN := wSMsgNewKey{
					Status: msg.Status,
					Type:   NewKey,
					Key:    key,
				}
				// sends all admins about key updated
				AddAllAdminMessage(stN)
			}
		}
	}
}

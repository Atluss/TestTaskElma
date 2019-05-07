// Package wsServer connection for admins
package wsServer

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/pkg/v1/dataKeys"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

var admins = make(map[*admin]bool) // admins connects

// handleNewKeyBroadcast send CPU load all accepted connections
func AddAllAdminMessage(msg interface{}) {
	for admin := range admins {
		admin.send <- msg
	}
}

type admin struct {
	Conn *websocket.Conn
	send chan interface{}
	run  chan string
}

// Run admin handlers
func (obj *admin) Run(db *gorm.DB) {
	admins[obj] = true
	go obj.Read(db)
	go obj.Write()
	wait := <-obj.run
	log.Println(wait)
}

// Stop connection and remove admin from admins list
func (obj *admin) Stop(err error) {
	if obj.Conn != nil {
		obj.Conn.Close()
	}
	delete(admins, obj)
	log.Printf("admins: len: %d %+v", len(admins), admins)
	obj.run <- fmt.Sprintf("stop: err %s", err)
}

// Write send messages
func (obj *admin) Write() {
	for msg := range obj.send {
		err := obj.Conn.WriteJSON(msg)
		if err != nil {
			obj.Stop(err)
			break
		}
	}
}

// Read process input messages
func (obj *admin) Read(db *gorm.DB) {
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
			Items:  []dataKeys.Keys{},
		}

		switch req.Type {
		case GetList:
			{
				msg.Type = GetList
				if req.Status > 10 { // for online clients
					msg.Items = GetActiveClients()
				} else {
					msg.Items, err = dataKeys.GetKeysByStatus(req.Status, db)
					if err != nil {
						msg.Status = http.StatusInternalServerError
					}
				}
				obj.send <- msg
			}
		case UpdateKey:
			{
				msg.Type = UpdateKey
				key := dataKeys.Keys{
					Key:    req.Key,
					Name:   req.Name,
					Status: req.Status,
				}

				if err := key.Update(db); err != nil {
					msg.Status = http.StatusInternalServerError
				} else {
					msg.UpdatedKey = key.Key
					if key.Status == dataKeys.KeyBlocked {
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

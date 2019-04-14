// websocket connection for clients
package ws_server

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib/data"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

type BannedKeys struct {
	mt sync.RWMutex
	V  map[string]bool
}

func (obj *BannedKeys) CloneMe() map[string]bool {
	obj.mt.RLock()
	nm := obj.V
	obj.mt.RUnlock()
	return nm
}

var BanKeys = BannedKeys{mt: sync.RWMutex{}, V: map[string]bool{}}
var Clients = make(map[*Client]bool) // connected clients

type Client struct {
	Conn *websocket.Conn
	Key  data.Keys
	run  chan string
}

func GetActiveClients() (keys []data.Keys) {
	for client := range Clients {
		keys = append(keys, client.Key)
	}
	return keys
}

func (obj *Client) Run() {
	obj.SendAdminsAboutMe(data.KeyOnline)
	for {
		st := data.Keys{}
		err := obj.Conn.ReadJSON(&st)
		if err != nil {
			log.Println("stop from run")
			obj.Stop(err)
			break
		}
		log.Printf("%+v", st)
	}
}

func (obj *Client) Stop(err error) {
	log.Printf("error: %+v", err)
	obj.Conn.Close()
	obj.SendAdminsAboutMe(data.KeyOfline)
	delete(Clients, obj)
	log.Printf("Clients len: %d", len(Clients))
	obj.run <- "Stop"
}

func (obj *Client) SendAdminsAboutMe(status int) {
	stN := wSMsgNewKey{
		Status: http.StatusOK,
		Type:   NewKey,
		Key: data.Keys{
			Key:    obj.Key.Key,
			Name:   obj.Key.Name,
			Ip:     obj.Key.Ip,
			Status: status,
		},
	}
	AddAllAdminMessage(stN)
}

func (obj *Client) Validate(db *gorm.DB) bool {

	err := obj.Conn.ReadJSON(&obj.Key)
	if err != nil {
		obj.Stop(err)
		return false
	}

	if obj.Key.Key == "" {
		obj.Stop(fmt.Errorf("error: client not send Key"))
		return false
	}

	if err := obj.Key.LoadByKey(db); err != nil {
		if err := obj.Key.Create(db); err != nil {
			log.Println(err)
			// send 500 if we can't create
			obj.WriteClose(websocket.CloseInternalServerErr, http.StatusText(http.StatusInternalServerError))
			return false
		} else {
			obj.SendAdminsAboutMe(obj.Key.Status)
		}
	}

	if obj.Key.Status == data.KeyBlocked {
		obj.WriteClose(websocket.ClosePolicyViolation, "banned")
		log.Printf("client Key: %s banned", obj.Key.Key)
		return false
	}

	if obj.Key.Status == data.KeyActive {
		// reg client for broadcast
		Clients[obj] = true
	} else {
		obj.WriteClose(websocket.CloseNormalClosure, "Key accepted, please wait and connected again.")
		return false
	}

	return true
}

func (obj *Client) WriteClose(clNum int, mgs string) {
	if err := obj.Conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(clNum, mgs),
		time.Time{}); err != nil {
		log.Printf("Can't close connection: %s", mgs)
	}
}

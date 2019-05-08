// Package wsServer websocket connection for clients
package wsServer

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/pkg/v1/dataKeys"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

// BannedKeys list banned keys when server run, this list generate when admin ban key
type BannedKeys struct {
	mt sync.RWMutex
	V  map[string]bool
}

// CloneMe return list banned keys safely
func (obj *BannedKeys) CloneMe() map[string]bool {
	obj.mt.RLock()
	nm := obj.V
	obj.mt.RUnlock()
	return nm
}

// BanKeys initialize list of banned keys
var BanKeys = BannedKeys{mt: sync.RWMutex{}, V: map[string]bool{}}

// Clients connected clients
var Clients = make(map[*Client]bool)

// Client typical client online :)
type Client struct {
	Conn *websocket.Conn
	Key  dataKeys.Keys
	run  chan string
}

// GetActiveClients return online keys right now
func GetActiveClients() (keys []dataKeys.Keys) {
	for client := range Clients {
		keys = append(keys, client.Key)
	}
	return keys
}

// Run process messages from client
func (obj *Client) Run() {
	obj.SendAdminsAboutMe(dataKeys.KeyOnline)
	for {
		st := dataKeys.Keys{}
		err := obj.Conn.ReadJSON(&st)
		if err != nil {
			log.Println("stop from run")
			obj.Stop(err)
			break
		}
		log.Printf("%+v", st)
	}
}

// Stop client connection and sends admin client off
func (obj *Client) Stop(err error) {
	log.Printf("error: %+v", err)
	obj.Conn.Close()
	obj.SendAdminsAboutMe(dataKeys.KeyOffline)
	delete(Clients, obj)
	log.Printf("Clients len: %d", len(Clients))
	obj.run <- "Stop"
}

// SendAdminsAboutMe sends all admins messages
func (obj *Client) SendAdminsAboutMe(status int) {
	stN := wSMsgNewKey{
		Status: http.StatusOK,
		Type:   NewKey,
		Key: dataKeys.Keys{
			Key:    obj.Key.Key,
			Name:   obj.Key.Name,
			Ip:     obj.Key.Ip,
			Status: status,
		},
	}
	AddAllAdminMessage(stN)
}

// Validate client connection and try accept it for broadcast
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
			// send 500 if server can't create
			obj.WriteClose(websocket.CloseInternalServerErr, http.StatusText(http.StatusInternalServerError))
			return false
		}
		obj.SendAdminsAboutMe(obj.Key.Status)
	}

	if obj.Key.Status == dataKeys.KeyBlocked {
		obj.WriteClose(websocket.ClosePolicyViolation, "banned")
		log.Printf("client Key: %s banned", obj.Key.Key)
		return false
	}

	if obj.Key.Status == dataKeys.KeyActive {
		// reg client for broadcast
		Clients[obj] = true
	} else {
		obj.WriteClose(websocket.CloseNormalClosure, "Key accepted, please wait and connected again.")
		return false
	}

	return true
}

// WriteClose send status close and message
func (obj *Client) WriteClose(clNum int, mgs string) {
	if err := obj.Conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(clNum, mgs),
		time.Time{}); err != nil {
		log.Printf("Can't close connection: %s", mgs)
	}
}

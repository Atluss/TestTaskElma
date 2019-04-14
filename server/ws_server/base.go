package ws_server

import (
	"github.com/Atluss/TestTaskElma/lib/data"
	"github.com/gorilla/websocket"
	"net/http"
)

type wSMsgNewKey struct {
	Status int
	Type   string
	Key    data.Keys
}

const (
	GetList   = "getList"
	UpdateKey = "updateKey"
	NewKey    = "newKey"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

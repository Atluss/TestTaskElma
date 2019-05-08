package wsServer

import (
	"github.com/Atluss/TestTaskElma/pkg/v1/dataKeys"
	"github.com/gorilla/websocket"
	"net/http"
)

type wSMsgNewKey struct {
	Status int
	Type   string
	Key    dataKeys.Keys
}

// types messages
const (
	GetList   = "getList"
	UpdateKey = "updateKey"
	NewKey    = "newKey"
)

// Upgrader ws config buffer
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // accept all origins
		return true
	},
}

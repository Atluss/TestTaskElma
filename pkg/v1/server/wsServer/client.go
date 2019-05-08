// Package wsServer clients websocket connection
package wsServer

import (
	"github.com/Atluss/TestTaskElma/pkg/v1/api"
	"github.com/Atluss/TestTaskElma/pkg/v1/config"
	"github.com/Atluss/TestTaskElma/pkg/v1/dataKeys"
	"log"
	"net/http"
)

// WSClient add websocket connection for clients
func WSClient(set *config.Setup) error {
	wsClient := &wSClient{
		Url:   "/st_cpu",
		Setup: set,
	}
	set.Route.HandleFunc(wsClient.Url, wsClient.HandleConnection)
	return nil
}

type wSClient struct {
	api.RequestRest
	Setup *config.Setup
	Url   string
}

// HandleConnection client connection
func (obj *wSClient) HandleConnection(w http.ResponseWriter, r *http.Request) {

	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	client := Client{
		Conn: ws,
		Key: dataKeys.Keys{
			Name: r.UserAgent(),
			Ip:   r.RemoteAddr,
		},
	}
	if client.Validate(obj.Setup.Gorm) {
		client.Run()
	}
}

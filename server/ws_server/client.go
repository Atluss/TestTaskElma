package ws_server

import (
	"github.com/Atluss/TestTaskElma/lib/api"
	"github.com/Atluss/TestTaskElma/lib/config"
	"github.com/Atluss/TestTaskElma/lib/data"
	"log"
	"net/http"
)

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

func (obj *wSClient) HandleConnection(w http.ResponseWriter, r *http.Request) {

	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	client := Client{
		Conn: ws,
		Key: data.Keys{
			Name: r.UserAgent(),
			Ip:   r.RemoteAddr,
		},
	}
	if client.Validate(obj.Setup.Gorm) {
		client.Run()
	}
}

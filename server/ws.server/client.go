package ws_server

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib/api"
	"github.com/Atluss/TestTaskElma/lib/config"
	"github.com/Atluss/TestTaskElma/lib/data"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func WSClient(set *config.Setup) error {

	wsClient := &wSClient{
		Url:   fmt.Sprintf("/%s/login", api.V1Api),
		Setup: set,
	}

	set.Route.HandleFunc("/st_cpu", wsClient.HandleConnection)

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
		log.Fatal(err)
	}
	defer ws.Close()

	st := data.Keys{}
	err = ws.ReadJSON(&st)
	if err != nil {
		log.Printf("error: %v", err)
		ws.CloseHandler()
		delete(Clients, ws)
		return
	}

	if st.Key == "" {
		log.Printf("error: client not send key")
		ws.CloseHandler()
		return
	}
	st.Name = r.UserAgent()
	st.Ip = r.RemoteAddr

	if err := st.LoadByKey(obj.Setup.Gorm); err != nil {
		if err := st.Create(obj.Setup.Gorm); err != nil {

			log.Println(err)
			// send 500 if we can't create
			if err := ws.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseInternalServerErr, http.StatusText(http.StatusInternalServerError)),
				time.Time{}); err != nil {
				log.Println("Can't close ws connection, by 500")
			}
			return
		}
	}

	if st.Status == 2 {
		if err := ws.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "banned"),
			time.Time{}); err != nil {
			log.Println("Can't close connection, by banned key")
		}
		log.Printf("client key: %s banned", st.Key)
		return
	}

	if st.Status == 1 {
		Clients[ws] = true
	} else {
		if err := ws.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Key accepted, please wait and connected again."),
			time.Time{}); err != nil {
			log.Println("Can't close connection: key accepted")
		}
		return
	}

	for {

		st := data.Keys{}
		err = ws.ReadJSON(&st)
		if err != nil {
			log.Printf("error: %v", err)
			delete(Clients, ws)
			break
		}

		log.Printf("%+v", st)
	}

}

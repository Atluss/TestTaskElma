package main

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib"
	"github.com/Atluss/TestTaskElma/lib/config"
	cpu "github.com/Atluss/TestTaskElma/lib/cpu.status"
	"github.com/Atluss/TestTaskElma/server/rest.api/v1"
	webserve "github.com/Atluss/TestTaskElma/server/web.server"
	ws_server "github.com/Atluss/TestTaskElma/server/ws.server"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Broadcast = make(chan cpu.CPULoad) // broadcast channel

func main() {

	path := "settings.json"
	set := config.NewApiSetup(path)
	set.Config.Print()

	// files for web pages: images, css, fonts, js and etc.
	set.Route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./http.files/public"))))

	// web pages
	webserve.AddPage("http.files/list.html", "/", true, set)
	webserve.AddPage("http.files/client.html", "/client", false, set)

	// setup rest endpoints
	lib.FailOnError(v1.V1Login(set, false), "error")
	lib.FailOnError(v1.V1Logout(set), "error")
	lib.FailOnError(v1.V1AddKey(set), "error")

	go cpu.GetCpuLoad(Broadcast)
	go HandleMessages()

	lib.FailOnError(ws_server.WSClient(set), "error")

	/*key := data.Keys{
		Key: "1231231232123",
	}

	if err := key.Create(set.Gorm); err != nil {
		log.Println(err)
	}*/

	/*set.Gorm.Create(&data.Keys{
		Key: "12312---454655644-3",
		Status: 0,
	})

	product := data.Keys{}
	set.Gorm.First(&product, "key = ?", "12312----3")*/

	lib.FailOnError(http.ListenAndServe(fmt.Sprintf(":%s", set.Config.Port), set.Route), "error")

}

func HandleMessages() {
	for {

		msg := <-Broadcast

		for client := range ws_server.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(ws_server.Clients, client)
			}
		}
	}
}

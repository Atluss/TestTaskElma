package main

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib"
	"github.com/Atluss/TestTaskElma/lib/config"
	cpu "github.com/Atluss/TestTaskElma/lib/cpu.status"
	"github.com/Atluss/TestTaskElma/server/rest.api/v1"
	webserve "github.com/Atluss/TestTaskElma/server/web.serve"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var Clients = make(map[*websocket.Conn]bool) // connected clients
var Broadcast = make(chan cpu.CPULoad)       // broadcast channel

func main() {

	path := "settings.json"
	set := config.NewApiSetup(path)
	set.Config.Print()

	// files for web pages images, css, fonts, js and etc.
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

	set.Route.HandleFunc("/st_cpu", HandleConnections)

	/*set.Gorm.AutoMigrate(&data.Keys{})
	set.Gorm.Create(&data.Keys{
		Key: "12312---454655644-3",
		Status: 0,
		Info: data.Info{
			Name: "a",
			Ip: "192.168.0.1",
		},
	})

	product := data.Keys{}
	set.Gorm.First(&product, "key = ?", "12312----3")
	log.Printf("%+v", product.Info.Ip)*/

	lib.FailOnError(http.ListenAndServe(fmt.Sprintf(":%s", set.Config.Port), set.Route), "error")

}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	Clients[ws] = true

	for {
		var msg cpu.CPULoad
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(Clients, ws)
			break
		}
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-Broadcast
		// Send it out to every client that is currently connected
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}

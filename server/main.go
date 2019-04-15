// This is execute file for server
// Run:
// 1. Web server
// 2. RestAPI
// 3. Websocket connection
// 4. Two go rutines for CPU broadcasts
package main

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib"
	"github.com/Atluss/TestTaskElma/lib/config"
	cpu "github.com/Atluss/TestTaskElma/lib/cpu_status"
	"github.com/Atluss/TestTaskElma/server/rest_api/v1"
	webserver "github.com/Atluss/TestTaskElma/server/web_server"
	"github.com/Atluss/TestTaskElma/server/ws_server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	path := "settings.json"
	set := config.NewApiSetup(path)
	set.Config.Print()

	// stop server
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Try stop server...")
		lib.LogOnError(set.Gorm.Close(), "error stop gorm")
		time.Sleep(time.Second)
		log.Println("Server stop! Buy buy!")
		os.Exit(1)
	}()

	// files for web pages: images, css, fonts, js and etc.
	set.Route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./http_files/public"))))

	// web pages
	webserver.AddPage("http_files/list.html", "/", true, set)
	webserver.AddPage("http_files/client.html", "/client", false, set)

	// setup rest endpoints
	lib.FailOnError(v1.V1Login(set, false), "error")
	lib.FailOnError(v1.V1Logout(set), "error")

	// broadcast cpu load two gorutines in there
	cpu.RunCPUBroadcast()

	//ws connection setup
	lib.FailOnError(ws_server.WSClient(set), "error")
	lib.FailOnError(ws_server.WSList(set, true), "error")

	lib.FailOnError(http.ListenAndServe(fmt.Sprintf(":%s", set.Config.Port), set.Route), "error")
}

package main

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib"
	"github.com/Atluss/TestTaskElma/lib/config"
	cpu "github.com/Atluss/TestTaskElma/lib/cpu_status"
	"github.com/Atluss/TestTaskElma/server/rest_api/v1"
	webserver "github.com/Atluss/TestTaskElma/server/web_server"
	"github.com/Atluss/TestTaskElma/server/ws_server"
	"net/http"
)

func main() {

	path := "settings.json"
	set := config.NewApiSetup(path)
	set.Config.Print()

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

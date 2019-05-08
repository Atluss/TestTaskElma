// This is execute file for server
// Run:
// 1. Web server
// 2. RestAPI
// 3. Websocket connection
// 4. Two go rutines for CPU broadcasts
package main

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/pkg/v1"
	"github.com/Atluss/TestTaskElma/pkg/v1/config"
	"github.com/Atluss/TestTaskElma/pkg/v1/cpu_status"
	"github.com/Atluss/TestTaskElma/pkg/v1/server/restApi"
	"github.com/Atluss/TestTaskElma/pkg/v1/server/webServer"
	"github.com/Atluss/TestTaskElma/pkg/v1/server/wsServer"
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
		v1.LogOnError(set.Gorm.Close(), "error stop gorm")
		time.Sleep(time.Second)
		log.Println("Server stop! Buy buy!")
		os.Exit(1)
	}()

	// files for web pages: images, css, fonts, js and etc.
	set.Route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./http_files/public"))))

	// web pages
	webServer.AddPage("http_files/list.html", "/", true, set)
	webServer.AddPage("http_files/client.html", "/client", false, set)

	// setup rest endpoints
	v1.FailOnError(restApi.V1Login(set, false), "error")
	v1.FailOnError(restApi.V1Logout(set), "error")

	// broadcast cpu load two gorutines in there
	cpuStatus.RunCPUBroadcast()

	//ws connection setup
	v1.FailOnError(wsServer.WSClient(set), "error")
	v1.FailOnError(wsServer.WSList(set, true), "error")

	v1.FailOnError(http.ListenAndServe(fmt.Sprintf(":%s", set.Config.Port), set.Route), "error")
}

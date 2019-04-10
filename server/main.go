package main

import (
	"github.com/Atluss/TestTaskElma/lib/config"
	webserve "github.com/Atluss/TestTaskElma/server/web.serve"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key         = []byte("super-secret-key")
	store       = sessions.NewCookieStore(key)
	sessionName = "session_a"
)

func main() {

	path := "settings.json"
	set := config.NewApiSetup(path)

	// files for web pages images, css, fonts, js and etc.
	set.Route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./http.files/public"))))

	webserve.AddPage("http.files/list.html", "/", true, set)

	//test login logout
	set.Route.HandleFunc("/login", login)
	set.Route.HandleFunc("/logout", logout)

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

	http.ListenAndServe(":8080", set.Route)

}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

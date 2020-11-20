package main

import (
	"log"
	"metafora-bot-service/controller"
	"metafora-bot-service/database"
	// "metafora-bot-service/googledrive"
	"metafora-bot-service/restfull"
	"metafora-bot-service/webapp"
	"net/http"
)

func main() {
	// srv := googledrive.Auth()
	db := database.New()
	nats := controller.Init(db)
	conn := webapp.NewServerSoket("/v1", nats/*, srv*/)
	go conn.ListenSocket()
	restfull.NewAdminInterface()
	restfull.NewUploadInterface()
	restfull.Init(db.GetDb, nats.GetTableName, "KJjvcj4545sd#jssdf7&sdf", true)
	http.HandleFunc("/upload", conn.UploadFile)
	log.Fatal(http.ListenAndServe(conn.Controller.Host, nil))
}

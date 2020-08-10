package main

import (
	"net/http"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio-demo/app/controllers"
	"github.com/kooinam/fabio-demo/app/models"
)

func main() {
	fab.Setup()

	models.RoomsCollection.Create()
	models.RoomsCollection.Create()
	models.RoomsCollection.Create()

	fab.RegisterController("session", &controllers.SessionsController{})
	fab.RegisterController("player", &controllers.PlayersController{})
	fab.RegisterController("room", &controllers.RoomsController{})

	fab.Serve(func() {
		fs := http.FileServer(http.Dir("./demo"))
		http.Handle("/demo/", http.StripPrefix("/demo/", fs))
	})
}

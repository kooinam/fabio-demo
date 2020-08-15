package main

import (
	"net/http"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio-demo/app/controllers"
	"github.com/kooinam/fabio-demo/app/models"
)

func main() {
	fab.Setup()

	models.PlayersCollection = fab.ModelManager().CreateCollection("player", models.MakePlayer)
	models.RoomsCollection = fab.ModelManager().CreateCollection("room", models.MakeRoom)

	models.RoomsCollection.Create()
	models.RoomsCollection.Create()
	models.RoomsCollection.Create()

	fab.ControllerManager().RegisterController("session", &controllers.SessionsController{})
	fab.ControllerManager().RegisterController("player", &controllers.PlayersController{})
	fab.ControllerManager().RegisterController("room", &controllers.RoomsController{})

	fab.ControllerManager().Serve(func() {
		fs := http.FileServer(http.Dir("./demo"))
		http.Handle("/demo/", http.StripPrefix("/demo/", fs))
	})
}

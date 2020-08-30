package main

import (
	"fmt"
	"net/http"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio-demo/app/controllers"
	"github.com/kooinam/fabio-demo/app/models"
	"github.com/kooinam/fabio/helpers"
	"github.com/kooinam/fabio/mongorecords"
	"github.com/kooinam/fabio/simplerecords"
)

func main() {
	fab.Setup()

	simpleAdapter := simplerecords.MakeAdapter()
	fab.ModelManager().RegisterAdapter("simple", simpleAdapter)

	mongoAdapter, err := mongorecords.MakeAdapter("mongodb://localhost:27017", "tic-tac-toe")
	if err != nil {
		panic("failed to connect to mongo")
	}
	fab.ModelManager().RegisterAdapter("mongo", mongoAdapter)

	models.PlayersCollection = fab.ModelManager().RegisterCollection("mongo", "players", models.MakePlayer)
	models.RoomsCollection = fab.ModelManager().RegisterCollection("mongo", "rooms", models.MakeRoom)

	roomNames := []string{"Room 1", "Room 2", "Room 3"}
	for _, roomName := range roomNames {
		result := models.RoomsCollection.Query().Where(helpers.H{
			"name": roomName,
		}).FirstOrCreate(helpers.H{})

		if !result.StatusSuccess() {
			panic(fmt.Sprintf("failed to load room %v", roomName))
		}

		result.Item().Memoize()
	}

	fab.ControllerManager().RegisterController("sessions", &controllers.SessionsController{})
	fab.ControllerManager().RegisterController("players", &controllers.PlayersController{})
	fab.ControllerManager().RegisterController("rooms", &controllers.RoomsController{})

	fab.ControllerManager().Serve("8000", func() {
		fs := http.FileServer(http.Dir("./demo"))
		http.Handle("/demo/", http.StripPrefix("/demo/", fs))
	})
}

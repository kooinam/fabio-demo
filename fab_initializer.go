package main

import (
	"fmt"
	"net/http"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio/controllers"
	"github.com/kooinam/fabio/helpers"
	"github.com/kooinam/fabio/models"
	"github.com/kooinam/fabio/mongorecords"
	"github.com/kooinam/fabio/simplerecords"
	"github.com/kooinam/fabio/views"

	Controllers "github.com/kooinam/fabio-demo/app/controllers"
	Models "github.com/kooinam/fabio-demo/app/models"
)

// FabInitializer initializes environment for fab.io usages
type FabInitializer struct{}

// Configure used to configure configurations
func (initializer *FabInitializer) Configure(configuration *fab.Configuration) {
	configuration.SetPort("8000")
	configuration.SetHttpHandler(func() {
		fs := http.FileServer(http.Dir("./demo"))
		http.Handle("/demo/", http.StripPrefix("/demo/", fs))
	})
}

// RegisterAdapters used to register adapters
func (initializer *FabInitializer) RegisterAdapters(manager *models.Manager) {
	mongoAdapter, err := mongorecords.MakeAdapter("mongodb://localhost:27017", "tic-tac-toe")
	if err != nil {
		panic("failed to connect to mongo")
	}
	manager.RegisterAdapter("mongo", mongoAdapter)

	simpleAdapter := simplerecords.MakeAdapter()
	manager.RegisterAdapter("simple", simpleAdapter)
}

// RegisterCollections used to register collections
func (initializer *FabInitializer) RegisterCollections(manager *models.Manager) {
	manager.RegisterCollection("mongo", "players", Models.MakePlayer)
	manager.RegisterCollection("mongo", "rooms", Models.MakeRoom)

	manager.RegisterCollection("simple", "bots", Models.MakeBot)
}

// RegisterControllers used to register controllers
func (initializer *FabInitializer) RegisterControllers(manager *controllers.Manager) {
	fab.ControllerManager().RegisterController("sessions", &Controllers.SessionsController{})
	fab.ControllerManager().RegisterController("players", &Controllers.PlayersController{})
	fab.ControllerManager().RegisterController("rooms", &Controllers.RoomsController{})
}

// RegisterViews used to register views
func (initializer *FabInitializer) RegisterViews(manager *views.Manager) {

}

// BeforeServe used for custom initializers after fab.io initializes and before serving
func (initializer *FabInitializer) BeforeServe() {
	// load rooms
	roomNames := []string{"Room 1", "Room 2", "Room 3", "Room 4", "Room 5", "Room 6"}

	for _, roomName := range roomNames {
		result := Models.RoomsCollection().Query().Where(helpers.H{
			"name": roomName,
		}).FirstOrCreate(helpers.H{})

		if !result.StatusSuccess() {
			panic(fmt.Sprintf("failed to load room %v", roomName))
		}

		result.Item().Memoize()
	}

	// load bots
	helpers.Times(10, func(i int) bool {
		result := Models.BotsCollection().Create(helpers.H{})

		if !result.StatusSuccess() {
			panic(fmt.Sprintf("faild to load bot"))
		}

		result.Item().Memoize()

		return true
	})
}

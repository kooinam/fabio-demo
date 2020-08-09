package main

import (
	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio-demo/app/controllers"
)

func main() {
	fab.Setup()

	fab.RegisterController("session", &controllers.SessionsController{})
	fab.RegisterController("player", &controllers.PlayersController{})
	fab.RegisterController("room", &controllers.RoomsController{})

	fab.Serve()
}

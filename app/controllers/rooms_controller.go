package controllers

import (
	"fmt"

	"github.com/kooinam/fabio-demo/app/models"

	"github.com/kooinam/fabio/controllers"
	Models "github.com/kooinam/fabio/models"
)

// RoomsController used for rooms actions
type RoomsController struct {
}

// AddBeforeActions used to add before actions callbacks
func (controller *RoomsController) AddBeforeActions(callbacksHandler *controllers.CallbacksHandler) {
	callbacksHandler.AddBeforeAction(controller.setCurrentPlayer)
	callbacksHandler.AddBeforeAction(controller.setCurrentRoom)
}

// AddActions used to add actions
func (controller *RoomsController) AddActions(actionsHandler *controllers.ActionsHandler) {
	actionsHandler.AddAction("List", controller.list)
	actionsHandler.AddAction("Join", controller.join)
	actionsHandler.AddAction("GrabSeat", controller.grabSeat)
	actionsHandler.AddAction("Leave", controller.leave)
	actionsHandler.AddAction("MakeMove", controller.makeMove)
}

func (controller *RoomsController) setCurrentPlayer(action string, connection *controllers.Connection) error {
	var err error

	authenticationToken := connection.Params("authenticationToken")

	currentPlayer := models.PlayersCollection.Find(func(base Models.Base) bool {
		player := base.(*models.Player)

		return player.GetAuthenticationToken() == authenticationToken
	})

	if currentPlayer != nil {
		connection.SetProperty("CurrentPlayer", currentPlayer)
	} else {
		err = fmt.Errorf("Unauthorized: %v", authenticationToken)
	}

	return err
}

// setCurrentRoom used to set current room
func (controller *RoomsController) setCurrentRoom(action string, connection *controllers.Connection) error {
	var err error

	roomID := connection.ParamsWithFallback("roomID", "").(string)

	currentRoom := models.RoomsCollection.FindByID(roomID)

	if currentRoom != nil {
		connection.SetProperty("CurrentRoom", currentRoom)
	}

	return err
}

func (controller *RoomsController) list(connection *controllers.Connection) (interface{}, error) {
	var err error
	var roomsView interface{}

	return roomsView, err
}

// join used for player to join a room
func (controller *RoomsController) join(connection *controllers.Connection) (interface{}, error) {
	var err error
	var roomView interface{}

	currentPlayer := connection.Property("CurrentPlayer").(*models.Player)
	room := models.JoinRoom(currentPlayer)

	connection.Join(room.ID)

	roomView = models.MakeRoomView(room, true)

	return roomView, err
}

// GrabSeat used for player to grab seat in room
func (controller *RoomsController) grabSeat(connection *controllers.Connection) (interface{}, error) {
	var err error
	var roomView interface{}

	currentPlayer := connection.Property("CurrentPlayer").(*models.Player)
	currentRoom, asserted := connection.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err = fmt.Errorf("room not found")

		return roomView, err
	}

	position, asserted := connection.ParamsWithFallback("position", -1).(float64)

	if asserted == false {
		err = fmt.Errorf("position is invalid")

		return roomView, err
	}

	err = currentRoom.GrabSeat(currentPlayer, int(position))

	if err != nil {
		return roomView, err
	}

	roomView = models.MakeRoomView(currentRoom, true)

	return roomView, err
}

// Leave used for player to leave a room
func (controller *RoomsController) leave(connection *controllers.Connection) (interface{}, error) {
	var err error
	var roomView interface{}

	currentPlayer := connection.Property("CurrentPlayer").(*models.Player)
	currentRoom, asserted := connection.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err = fmt.Errorf("room not found")

		return roomView, err
	}

	err = currentRoom.Leave(currentPlayer)

	if err != nil {
		return roomView, err
	}

	roomView = models.MakeRoomView(currentRoom, true)

	return roomView, err
}

// MakeMove used for player to make move in room
func (controller *RoomsController) makeMove(connection *controllers.Connection) (interface{}, error) {
	var err error
	var roomView interface{}

	currentPlayer := connection.Property("CurrentPlayer").(*models.Player)
	currentRoom, asserted := connection.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err = fmt.Errorf("room not found")

		return roomView, err
	}

	x := int(connection.ParamsWithFallback("x", -1).(float64))
	y := int(connection.ParamsWithFallback("y", -1).(float64))

	err = currentRoom.MakeMove(currentPlayer, x, y)

	if err == nil {
		roomView = models.MakeRoomView(currentRoom, true)
	}

	return roomView, err
}

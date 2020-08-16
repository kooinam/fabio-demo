package controllers

import (
	"fmt"

	fab "github.com/kooinam/fabio"

	"github.com/kooinam/fabio-demo/app/models"

	"github.com/kooinam/fabio/controllers"
	"github.com/kooinam/fabio/helpers"
	Models "github.com/kooinam/fabio/models"
)

// RoomsController is controller for room's actions
type RoomsController struct {
}

// RegisterBeforeHooks used to register before action hooks
func (controller *RoomsController) RegisterBeforeHooks(hooksHandler *controllers.HooksHandler) {
	hooksHandler.RegisterBeforeHook(controller.setCurrentPlayer)
	hooksHandler.RegisterBeforeHook(controller.setCurrentRoom)
}

// RegisterActions used to register actions
func (controller *RoomsController) RegisterActions(actionsHandler *controllers.ActionsHandler) {
	actionsHandler.RegisterAction("List", controller.list)
	actionsHandler.RegisterAction("Join", controller.join)
	actionsHandler.RegisterAction("GrabSeat", controller.grabSeat)
	actionsHandler.RegisterAction("Leave", controller.leave)
	actionsHandler.RegisterAction("MakeMove", controller.makeMove)
}

func (controller *RoomsController) setCurrentPlayer(action string, connection *controllers.Context) error {
	var err error

	authenticationToken := connection.ParamsStr("authenticationToken")

	currentPlayer := models.PlayersCollection.Find(func(item Models.Modellable) bool {
		player := item.(*models.Player)

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
func (controller *RoomsController) setCurrentRoom(action string, context *controllers.Context) error {
	var err error

	roomID := context.ParamsStr("roomID")

	currentRoom := models.RoomsCollection.FindByID(roomID)

	if currentRoom != nil {
		context.SetProperty("CurrentRoom", currentRoom)
	}

	return err
}

func (controller *RoomsController) list(context *controllers.Context) (interface{}, error) {
	var err error
	var roomsView interface{}

	rooms := models.AssertRooms(models.RoomsCollection.GetItems())
	roomsView = models.MakeRoomsView(rooms)

	return roomsView, err
}

// join used for player to join a room
func (controller *RoomsController) join(context *controllers.Context) (interface{}, error) {
	var err error
	var roomView interface{}

	currentPlayer := context.Property("CurrentPlayer").(*models.Player)
	roomID := context.ParamsStr("roomID")

	room, err := models.JoinRoom(currentPlayer, roomID)

	if err != nil {
		return roomView, err
	}

	context.SingleJoin(room.ID)

	roomView = models.MakeRoomView(room, true)

	return roomView, err
}

// GrabSeat used for player to grab seat in room
func (controller *RoomsController) grabSeat(context *controllers.Context) (interface{}, error) {
	var err error
	var roomView interface{}

	currentRoom, asserted := context.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err = fmt.Errorf("room not found")

		return roomView, err
	}

	currentPlayer := context.Property("CurrentPlayer").(*models.Player)
	position := context.Params("position")

	err = fab.ActorManager().Request(currentRoom.Actor.Identifier(), "GrabSeat", helpers.H{
		"playerID": currentPlayer.ID,
		"position": position,
	})

	if err != nil {
		return roomView, err
	}

	roomView = models.MakeRoomView(currentRoom, true)

	return roomView, err
}

// Leave used for player to leave a room
func (controller *RoomsController) leave(context *controllers.Context) (interface{}, error) {
	var err error
	var roomView interface{}

	currentRoom, asserted := context.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err = fmt.Errorf("room not found")

		return roomView, err
	}

	currentPlayer := context.Property("CurrentPlayer").(*models.Player)

	err = fab.ActorManager().Request(currentRoom.Actor.Identifier(), "Leave", helpers.H{
		"playerID": currentPlayer.ID,
	})

	if err != nil {
		return roomView, err
	}

	roomView = models.MakeRoomView(currentRoom, true)

	return roomView, err
}

// MakeMove used for player to make move in room
func (controller *RoomsController) makeMove(context *controllers.Context) (interface{}, error) {
	var err error
	var roomView interface{}

	currentRoom, asserted := context.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err = fmt.Errorf("room not found")

		return roomView, err
	}

	currentPlayer := context.Property("CurrentPlayer").(*models.Player)
	x := context.Params("x")
	y := context.Params("y")

	err = fab.ActorManager().Request(currentRoom.Actor.Identifier(), "MakeMove", helpers.H{
		"playerID": currentPlayer.ID,
		"x":        x,
		"y":        y,
	})

	if err == nil {
		roomView = models.MakeRoomView(currentRoom, true)
	}

	return roomView, err
}

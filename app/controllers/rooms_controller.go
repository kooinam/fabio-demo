package controllers

import (
	"fmt"

	fab "github.com/kooinam/fabio"

	"github.com/kooinam/fabio-demo/app/models"

	"github.com/kooinam/fabio/controllers"
	"github.com/kooinam/fabio/helpers"
)

// RoomsController is controller for room's actions
type RoomsController struct {
}

// RegisterHooksAndActions used to register hooks and actions
func (controller *RoomsController) RegisterHooksAndActions(hooksHandler *controllers.HooksHandler, actionsHandler *controllers.ActionsHandler) {
	actionsHandler.RegisterConnectedAction(controller.connected)

	hooksHandler.RegisterBeforeActionHook(controller.setCurrentPlayer)
	hooksHandler.RegisterBeforeActionHook(controller.setCurrentRoom)

	actionsHandler.RegisterAction("List", controller.list)
	actionsHandler.RegisterAction("Join", controller.join)
	actionsHandler.RegisterAction("GrabSeat", controller.grabSeat)
	actionsHandler.RegisterAction("Leave", controller.leave)
	actionsHandler.RegisterAction("MakeMove", controller.makeMove)
}

func (controller *RoomsController) connected(context *controllers.Context) {
	context.Join("lobby")
}

func (controller *RoomsController) setCurrentPlayer(action string, context *controllers.Context) {
	authenticationToken := context.ParamsStr("authenticationToken")

	result := models.FindPlayerByToken(authenticationToken)

	if !result.StatusSuccess() {
		// return aunthorize if player not found
		err := fmt.Errorf("Unauthorized: %v", authenticationToken)
		context.SetErrorResult(controllers.StatusUnauthorized, err)
		return
	}

	context.SetProperty("CurrentPlayer", result.Item())
}

// setCurrentRoom used to set current room
func (controller *RoomsController) setCurrentRoom(action string, context *controllers.Context) {
	roomID := context.ParamsStr("roomID")

	room := models.RoomsCollection().List().FindByID(roomID)

	if room != nil {
		context.SetProperty("CurrentRoom", room)
	}
}

func (controller *RoomsController) list(context *controllers.Context) {
	// results := models.RoomsCollection.Query().All()

	// if results.StatusError() {
	// 	context.SetErrorResult(controllers.StatusError, results.Error())
	// 	return
	// }

	// rooms := models.AssertRooms(results.Items())
	// roomsView := models.MakeRoomsView(rooms)

	// context.SetSuccessResult(roomsView)

	items := models.RoomsCollection().List().Items()

	rooms := models.AssertRooms(items)
	roomsView := models.MakeRoomsView(rooms)

	context.SetSuccessResult(roomsView)

}

// join used for player to join a room
func (controller *RoomsController) join(context *controllers.Context) {
	currentPlayer := context.Property("CurrentPlayer").(*models.Player)
	roomID := context.ParamsStr("roomID")

	room, err := models.JoinRoom(currentPlayer, roomID)

	if err != nil {
		context.SetErrorResult(controllers.StatusError, err)
		return
	}

	context.SingleJoin(room.GetID())

	roomView := models.MakeRoomView(room, true)

	context.SetSuccessResult(roomView)
}

// GrabSeat used for player to grab seat in room
func (controller *RoomsController) grabSeat(context *controllers.Context) {
	currentRoom, asserted := context.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err := fmt.Errorf("room not found")
		context.SetErrorResult(controllers.StatusError, err)
		return
	}

	currentPlayer := context.Property("CurrentPlayer").(*models.Player)
	position := context.Params("position")

	err := fab.ActorManager().Request(currentRoom.Actor().Identifier(), "GrabSeat", helpers.H{
		"playerID": currentPlayer.GetID(),
		"position": position,
	})

	if err != nil {
		context.SetErrorResult(controllers.StatusError, err)
		return
	}

	roomView := models.MakeRoomView(currentRoom, true)

	context.SetSuccessResult(roomView)
}

// Leave used for player to leave a room
func (controller *RoomsController) leave(context *controllers.Context) {
	currentRoom, asserted := context.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err := fmt.Errorf("room not found")
		context.SetErrorResult(controllers.StatusError, err)
		return
	}

	currentPlayer := context.Property("CurrentPlayer").(*models.Player)

	err := fab.ActorManager().Request(currentRoom.Actor().Identifier(), "Leave", helpers.H{
		"playerID": currentPlayer.GetID(),
	})

	if err != nil {
		context.SetErrorResult(controllers.StatusError, err)
		return
	}

	context.Leave(currentRoom.GetID())

	roomView := models.MakeRoomView(currentRoom, true)

	context.SetSuccessResult(roomView)
}

// MakeMove used for player to make move in room
func (controller *RoomsController) makeMove(context *controllers.Context) {
	currentRoom, asserted := context.Property("CurrentRoom").(*models.Room)

	if asserted == false {
		err := fmt.Errorf("room not found")
		context.SetErrorResult(controllers.StatusError, err)
		return
	}

	currentPlayer := context.Property("CurrentPlayer").(*models.Player)
	x := context.Params("x")
	y := context.Params("y")

	err := fab.ActorManager().Request(currentRoom.Actor().Identifier(), "MakeMove", helpers.H{
		"playerID": currentPlayer.GetID(),
		"x":        x,
		"y":        y,
	})

	if err != nil {
		context.SetErrorResult(controllers.StatusError, err)
		return
	}

	roomView := models.MakeRoomView(currentRoom, true)

	context.SetSuccessResult(roomView)
}

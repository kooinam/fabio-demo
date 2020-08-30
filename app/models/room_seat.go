package models

import (
	"fmt"

	"github.com/kooinam/fabio/helpers"

	"github.com/kooinam/fabio/models"
	"github.com/kooinam/fabio/simplerecords"
)

// RoomSeat is used to represent seat in room
type RoomSeat struct {
	simplerecords.Base
	Position    int `json:"position"`
	player      *Player
	hasMadeMove bool
}

// assertRoomSeats used to convert items to room seats
func assertRoomSeats(items []models.Modellable) []*RoomSeat {
	seats := make([]*RoomSeat, len(items))

	for i, item := range items {
		seats[i] = item.(*RoomSeat)
	}

	return seats
}

func makeRoomSeat(collection *models.Collection, hooksHandler *models.HooksHandler) models.Modellable {
	seat := &RoomSeat{}

	hooksHandler.RegisterInitializeHook(seat.initialize)

	return seat
}

func (seat *RoomSeat) initialize(attributes *helpers.Dictionary) {
	seat.Position = attributes.ValueInt("position", 0)
}

// Grab used to grab seat
func (seat *RoomSeat) Grab(player *Player) error {
	var err error

	if seat.isEmpty() == false {
		err = fmt.Errorf("seat is taken")

		return err
	}

	seat.player = player

	return err
}

// Leave used to leave seat
func (seat *RoomSeat) Leave(room *Room) error {
	var err error

	if room.State.Equals(RoomStates.Playing) {
		err = fmt.Errorf("cannot leave while playing")

		return err
	}

	if seat.isEmpty() == false {
		seat.player = nil
	}

	return err
}

func (seat *RoomSeat) isEmpty() bool {
	return seat.player == nil
}

func (seat *RoomSeat) isGrabbedBy(player *Player) bool {
	isGrabbedByPlayer := false

	if seat.player != nil && seat.player.GetID() == player.GetID() {
		isGrabbedByPlayer = true
	}

	return isGrabbedByPlayer
}

func (seat *RoomSeat) equals(aSeat *RoomSeat) bool {
	return seat.Position == aSeat.Position
}

func (seat *RoomSeat) reset() {
	seat.hasMadeMove = false
}

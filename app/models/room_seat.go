package models

import (
	"fmt"

	"github.com/kooinam/fabio/models"
)

// RoomSeat is used to represent seat in room
type RoomSeat struct {
	ID          string
	Position    int `json:"position"`
	player      *Player
	hasMadeMove bool
}

func makeRoomSeat(args ...interface{}) models.Base {
	return &RoomSeat{
		Position: args[0].(int),
	}
}

// assertRoomSeats used to convert items to room seats
func assertRoomSeats(items []models.Base) []*RoomSeat {
	seats := make([]*RoomSeat, len(items))

	for i, seat := range items {
		seats[i] = seat.(*RoomSeat)
	}

	return seats
}

// GetID used to get ID
func (seat *RoomSeat) GetID() string {
	return seat.ID
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

	if seat.player != nil && seat.player.ID == player.ID {
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

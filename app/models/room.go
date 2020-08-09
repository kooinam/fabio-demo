package models

import (
	"fmt"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio/helpers"
	"github.com/kooinam/fabio/models"
)

const gridSize = 3

// RoomsCollection is singleton for RoomsCollection
var RoomsCollection *models.Collection

func init() {
	RoomsCollection = models.MakeCollection(makeRoom)
}

// Room used to represent room data
type Room struct {
	ID       string
	runner   *models.Runner
	State    *models.FiniteStateMachine `json:"-"`
	Seats    *models.Collection         `json:"-"`
	Cells    [gridSize][gridSize]int    `json:"cells"`
	Rankings map[string]int             `json:"rankings"`
}

func makeRoom(args ...interface{}) models.Base {
	room := &Room{
		ID:       fmt.Sprintf("%v", RoomsCollection.Count()+1),
		Seats:    models.MakeCollection(makeRoomSeat),
		Rankings: make(map[string]int),
	}

	room.State = makeRoomFSM(room)
	room.runner = models.MakeRunner(room.Run, 2)

	helpers.Times(2, func(i int) bool {
		room.Seats.Create(i)

		return true
	})

	room.resetCells()

	room.runner.Ch <- 1 // start room runner

	return room
}

// GetID used to get ID
func (room *Room) GetID() string {
	return room.ID
}

// JoinRoom used to join a new room
func JoinRoom(player *Player) *Room {
	availableRoom := RoomsCollection.FindOrCreate(func(base models.Base) bool {
		return true
	}).(*Room)

	return availableRoom
}

// GrabSeat used to grab a seat in room
func (room *Room) GrabSeat(player *Player, position int) error {
	var err error

	originalSeat, asserted := room.Seats.Find(func(base models.Base) bool {
		seat := base.(*RoomSeat)

		return seat.isGrabbedBy(player)
	}).(*RoomSeat)

	if asserted {
		err = originalSeat.Leave(room)

		if err != nil {
			return err
		}
	}

	foundSeat, asserted := room.Seats.Find(func(base models.Base) bool {
		seat := base.(*RoomSeat)

		return seat.Position == position
	}).(*RoomSeat)

	if asserted == false {
		err = fmt.Errorf("seat is not found")

		if originalSeat != nil {
			// grab back original seat if got error
			originalSeat.Grab(player)
		}

		return err
	}

	err = foundSeat.Grab(player)

	if err != nil {
		if originalSeat != nil {
			// grab back original seat if got error
			originalSeat.Grab(player)
		}

		return err
	}

	roomView := MakeRoomView(room, true)
	parameters := fab.H{
		"seat": foundSeat.Position,
	}

	fab.BroadcastEvent("room", room.ID, "GrabbedSeat", roomView, parameters)

	return err
}

// Leave used to leave room
func (room *Room) Leave(player *Player) error {
	var err error

	originalSeat, asserted := room.Seats.Find(func(base models.Base) bool {
		seat := base.(*RoomSeat)

		return seat.isGrabbedBy(player)
	}).(*RoomSeat)

	if asserted {
		err = originalSeat.Leave(room)
	}

	if err != nil {
		return err
	}

	if originalSeat != nil {
		roomView := MakeRoomView(room, true)
		parameters := fab.H{
			"seat": originalSeat.Position,
		}

		fab.BroadcastEvent("room", room.ID, "LeftSeat", roomView, parameters)
	}

	return err
}

// MakeMove for player to make move in room
func (room *Room) MakeMove(player *Player, x int, y int) error {
	var err error

	if x < 0 || x >= gridSize || y < 0 || y >= gridSize {
		err = fmt.Errorf("index is invalid")

		return err
	}

	if room.State.Equals(RoomStates.Playing) == false {
		err = fmt.Errorf("not playing")

		return err
	}

	playerSeat, asserted := room.Seats.Find(func(base models.Base) bool {
		seat := base.(*RoomSeat)

		return seat.isGrabbedBy(player)
	}).(*RoomSeat)

	if asserted == false {
		err = fmt.Errorf("grab seat before making move")

		return err
	}

	activeSeat := room.State.GetActiveAgent().(*RoomSeat)

	if activeSeat.equals(playerSeat) == false {
		err = fmt.Errorf("player is not active")

		return err
	}

	if activeSeat.hasMadeMove {
		err = fmt.Errorf("player has made move")

		return err
	}

	cell := room.Cells[x][y]

	if cell != -1 {
		err = fmt.Errorf("cell is taken")

		return err
	}

	// cell is not taken. make move
	activeSeat.hasMadeMove = true
	room.Cells[x][y] = playerSeat.Position

	// broadcast make move event
	roomView := MakeRoomView(room, true)
	parameters := fab.H{
		"seat":  activeSeat.Position,
		"cellX": x,
		"cellY": y,
	}

	fab.BroadcastEvent("room", room.ID, "MadeMove", roomView, parameters)

	return err
}

// GetSeats used to get RoomSeat collection
func (room *Room) GetSeats() []*RoomSeat {
	seats := assertRoomSeats(room.Seats.GetItems())

	return seats
}

// Run to run room every regularly
func (room *Room) Run() {
	room.State.Run(room)
}

func (room *Room) resetCells() {
	for x := range room.Cells {
		for y := range room.Cells[x] {
			room.Cells[x][y] = -1
		}
	}
}

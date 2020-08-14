package models

import (
	"fmt"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio/actors"
	"github.com/kooinam/fabio/helpers"
	"github.com/kooinam/fabio/models"
)

const gridSize = 3

// RoomsCollection is singleton for RoomsCollection
var RoomsCollection *models.Collection

// Room used to represent room data
type Room struct {
	models.Base
	actor    *actors.Actor
	State    *models.FiniteStateMachine `json:"-"`
	Seats    *models.Collection         `json:"-"`
	Cells    [gridSize][gridSize]int    `json:"cells"`
	Rankings map[string]int             `json:"rankings"`
}

// MakeRoom used to instantiate room
func MakeRoom(collection *models.Collection, args ...interface{}) models.Modellable {
	room := &Room{}

	room.Initialize(collection)

	room.Seats = fab.ModelManager().CreateCollection("room_seats", makeRoomSeat)
	room.Rankings = make(map[string]int)

	room.actor = fab.ActorManager().RegisterActor(room)
	room.State = makeRoomFSM(room)

	helpers.Times(2, func(i int) bool {
		room.Seats.Create(i)

		return true
	})

	room.resetCells()

	return room
}

// AssertRooms used to convert items to rooms
func AssertRooms(items []models.Modellable) []*Room {
	rooms := make([]*Room, len(items))

	for i, item := range items {
		rooms[i] = item.(*Room)
	}

	return rooms
}

// JoinRoom used to join a new room
func JoinRoom(player *Player, roomID string) (*Room, error) {
	var err error

	availableRoom, asserted := RoomsCollection.FindByID(roomID).(*Room)

	if asserted == false {
		err = fmt.Errorf("room not found")

		return nil, err
	}

	return availableRoom, err
}

func (room *Room) RegisterActions(actionsHandler *actors.ActionsHandler) {
	actionsHandler.RegisterAction("Update", room.Run)
}

// GrabSeat used to grab a seat in room
func (room *Room) GrabSeat(player *Player, position int) error {
	var err error

	originalSeat, asserted := room.Seats.Find(func(item models.Modellable) bool {
		seat := item.(*RoomSeat)

		return seat.isGrabbedBy(player)
	}).(*RoomSeat)

	if asserted {
		err = originalSeat.Leave(room)

		if err != nil {
			return err
		}
	}

	foundSeat, asserted := room.Seats.Find(func(item models.Modellable) bool {
		seat := item.(*RoomSeat)

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

	fab.ControllerManager().BroadcastEvent("room", room.GetID(), "GrabbedSeat", roomView, parameters)

	return err
}

// Leave used to leave room
func (room *Room) Leave(player *Player) error {
	var err error

	originalSeat, asserted := room.Seats.Find(func(item models.Modellable) bool {
		seat := item.(*RoomSeat)

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

		fab.ControllerManager().BroadcastEvent("room", room.GetID(), "LeftSeat", roomView, parameters)
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

	playerSeat, asserted := room.Seats.Find(func(item models.Modellable) bool {
		seat := item.(*RoomSeat)

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

	fab.ControllerManager().BroadcastEvent("room", room.GetID(), "MadeMove", roomView, parameters)

	return err
}

// GetSeats used to get RoomSeat collection
func (room *Room) GetSeats() []*RoomSeat {
	seats := assertRoomSeats(room.Seats.GetItems())

	return seats
}

// Run to run room every regularly
func (room *Room) Run(context *actors.Context) error {
	var err error

	room.State.Run(room)

	return err
}

func (room *Room) resetCells() {
	for x := range room.Cells {
		for y := range room.Cells[x] {
			room.Cells[x][y] = -1
		}
	}
}

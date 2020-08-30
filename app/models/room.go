package models

import (
	"fmt"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio/actors"
	"github.com/kooinam/fabio/helpers"
	"github.com/kooinam/fabio/models"
	"github.com/kooinam/fabio/mongorecords"
)

const gridSize = 3

// RoomsCollection is singleton for RoomsCollection
var RoomsCollection *models.Collection

// Room used to represent room data
type Room struct {
	mongorecords.Base `bson:"base,inline"`
	Name              string                     `bson:"name" json:"name"`
	State             *models.FiniteStateMachine `bson:"-" json:"-"`
	Seats             *models.Collection         `json:"-"`
	Cells             [gridSize][gridSize]int    `json:"cells"`
	Rankings          map[string]int             `json:"rankings"`
	MembersCount      int                        `json:"membersCount"`
	MaxMembersCount   int                        `json:"maxMembersCount"`
	actor             *actors.Actor              `bson:"-"`
}

// AssertRooms used to convert items to rooms
func AssertRooms(items []models.Modellable) []*Room {
	rooms := make([]*Room, len(items))

	for i, item := range items {
		rooms[i] = item.(*Room)
	}

	return rooms
}

// MakeRoom used to instantiate room
func MakeRoom(collection *models.Collection, hooksHandler *models.HooksHandler) models.Modellable {
	room := &Room{}

	hooksHandler.RegisterInitializeHook(room.initialize)

	hooksHandler.RegisterAfterMemoizeHook(room.afterMemoize)

	return room
}

func (room *Room) initialize(attributes *helpers.Dictionary) {
	room.Name = attributes.ValueStr("name")
	room.State = makeRoomFSM(room)
	room.Rankings = make(map[string]int)

	room.Seats = fab.ModelManager().CreateCollection("simple", "room_seats", makeRoomSeat)

	helpers.Times(2, func(i int) bool {
		result := room.Seats.Create(helpers.H{
			"position": i,
		})

		if result.StatusSuccess() {
			result.Item().Memoize()
		}

		return true
	})

	room.resetCells()
}

func (room *Room) afterMemoize() {
	room.actor = fab.ActorManager().RegisterActor(room.GetCollectionName(), room)
}

// Actor used to retrieve room's actor
func (room *Room) Actor() *actors.Actor {
	if room.actor == nil {
		panic(fmt.Sprintf("no actor registered for room %v", room.Name))
	}
	return room.actor
}

// RegisterActorActions used to register actor's actions
func (room *Room) RegisterActorActions(actionsHandler *actors.ActionsHandler) {
	actionsHandler.RegisterAction("Update", room.run)
	actionsHandler.RegisterAction("GrabSeat", room.grabSeat)
	actionsHandler.RegisterAction("Leave", room.leave)
	actionsHandler.RegisterAction("MakeMove", room.makeMove)
}

// JoinRoom used to join a new room
func JoinRoom(player *Player, roomID string) (*Room, error) {
	// result := RoomsCollection.Query().Find(roomID)

	// if result.StatusError() {
	// 	return nil, result.Error()
	// } else if result.StatusNotFound() {
	// 	return nil, helpers.NotFoundError("room")
	// }

	// return result.Item().(*Room), nil

	var err error

	item := RoomsCollection.List().FindByID(roomID)

	if item == nil {
		err = fmt.Errorf("room not found")
	}

	return item.(*Room), err
}

// GetSeats used to get RoomSeat collection
func (room *Room) GetSeats() []*RoomSeat {
	seats := assertRoomSeats(room.Seats.List().Items())

	return seats
}

// run to run room regularly
func (room *Room) run(context *actors.Context) error {
	var err error

	room.State.Run(room)

	return err
}

// grabSeat used to grab a seat in room
func (room *Room) grabSeat(context *actors.Context) error {
	var err error

	playerID := context.ParamsStr("playerID")
	result := PlayersCollection.Query().Find(playerID)

	if result.StatusError() {
		return result.Error()
	} else if result.StatusNotFound() {
		return helpers.NotFoundError("player")
	}

	player := result.Item().(*Player)

	position := context.ParamsInt("position", -1)

	originalSeat, asserted := room.Seats.List().Find(func(item models.Modellable) bool {
		seat := item.(*RoomSeat)

		return seat.isGrabbedBy(player)
	}).(*RoomSeat)

	if asserted {
		err = originalSeat.Leave(room)

		if err != nil {
			return err
		}
	}

	foundSeat, asserted := room.Seats.List().Find(func(item models.Modellable) bool {
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
	parameters := helpers.H{
		"seat": foundSeat.Position,
	}

	fab.ControllerManager().BroadcastEvent(room.GetCollectionName(), room.GetID(), "GrabbedSeat", roomView, parameters)
	room.calculateMembers()

	return err
}

// leave used to leave room
func (room *Room) leave(context *actors.Context) error {
	var err error

	playerID := context.ParamsStr("playerID")
	result := PlayersCollection.Query().Find(playerID)

	if result.StatusError() {
		return result.Error()
	} else if result.StatusNotFound() {
		return helpers.NotFoundError("player")
	}

	player := result.Item().(*Player)

	originalSeat, asserted := room.Seats.List().Find(func(item models.Modellable) bool {
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
		parameters := helpers.H{
			"seat": originalSeat.Position,
		}

		fab.ControllerManager().BroadcastEvent(room.GetCollectionName(), room.GetID(), "LeftSeat", roomView, parameters)

		room.calculateMembers()
	}

	return err
}

// makeMove for player to make move in room
func (room *Room) makeMove(context *actors.Context) error {
	var err error

	x := context.ParamsInt("x", -1)
	y := context.ParamsInt("y", -1)
	playerID := context.ParamsStr("playerID")
	result := PlayersCollection.Query().Find(playerID)

	if result.StatusError() {
		return result.Error()
	} else if result.StatusNotFound() {
		return helpers.NotFoundError("player")
	}

	player := result.Item().(*Player)

	if x < 0 || x >= gridSize || y < 0 || y >= gridSize {
		err = fmt.Errorf("index is invalid")

		return err
	}

	if room.State.Equals(RoomStates.Playing) == false {
		err = fmt.Errorf("not playing")

		return err
	}

	playerSeat, asserted := room.Seats.List().Find(func(item models.Modellable) bool {
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

	if playerSeat.hasMadeMove {
		err = fmt.Errorf("player has made move")

		return err
	}

	cell := room.Cells[x][y]

	if cell != -1 {
		err = fmt.Errorf("cell is taken")

		return err
	}

	room.populateCell(playerSeat, x, y)

	return err
}

func (room *Room) populateCell(seat *RoomSeat, x int, y int) {
	// cell is not taken. make move
	seat.hasMadeMove = true
	room.Cells[x][y] = seat.Position

	// broadcast make move event
	roomView := MakeRoomView(room, true)
	parameters := helpers.H{
		"seat":  seat.Position,
		"cellX": x,
		"cellY": y,
	}

	fab.ControllerManager().BroadcastEvent(room.GetCollectionName(), room.GetID(), "MadeMove", roomView, parameters)
}

func (room *Room) resetCells() {
	for x := range room.Cells {
		for y := range room.Cells[x] {
			room.Cells[x][y] = -1
		}
	}
}

func (room *Room) calculateMembers() {
	grabbedSeats := room.Seats.List().FindAll(func(item models.Modellable) bool {
		seat := item.(*RoomSeat)

		return !seat.isEmpty()
	})

	room.MaxMembersCount = room.Seats.List().Count()
	room.MembersCount = len(grabbedSeats)

	roomView := MakeSimpleRoomView(room, true)

	fab.ControllerManager().BroadcastEvent(room.GetCollectionName(), "lobby", "RoomUpdated", roomView, helpers.H{})
}

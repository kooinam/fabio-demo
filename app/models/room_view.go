package models

import (
	"github.com/kooinam/fabio/helpers"
)

// RoomView used to represent room's view
type RoomView struct {
	*Room
	SeatViews          []*RoomSeatView `json:"seats"`
	State              string          `json:"state"`
	TurnEndAt          int64           `json:"turnEndAt"`
	ActiveSeatPosition int             `json:"activeSeatPosition"`
}

// MakeRoomView used to instantiate room view
func MakeRoomView(room *Room, includeRoot bool) interface{} {
	seatViews := make([]*RoomSeatView, room.Seats.Count())

	for i, seat := range room.GetSeats() {
		seatViews[i] = MakeRoomSeatView(seat)
	}

	activeSeat := room.State.GetActiveAgent()
	activeSeatPosition := -1

	if activeSeat != nil {
		activeSeatPosition = activeSeat.(*RoomSeat).Position
	}

	roomView := &RoomView{
		Room:               room,
		SeatViews:          seatViews,
		State:              room.State.GetName(),
		TurnEndAt:          room.State.GetEndAt(),
		ActiveSeatPosition: activeSeatPosition,
	}

	view := helpers.IncludeRootInJSON(roomView, includeRoot, "room")

	return view
}

func MakeRoomsView(rooms []*Room) interface{} {
	roomsView := make([]interface{}, len(rooms))
	for i, room := range rooms {
		roomsView[i] = MakeRoomView(room, false)
	}

	view := helpers.IncludeRootInJSON(roomsView, true, "rooms")

	return view
}

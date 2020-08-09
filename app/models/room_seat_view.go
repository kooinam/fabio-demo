package models

type RoomSeatView struct {
	*RoomSeat
	PlayerView interface{} `json:"player"`
}

func MakeRoomSeatView(roomSeat *RoomSeat) *RoomSeatView {
	roomSeatView := &RoomSeatView{
		RoomSeat:   roomSeat,
		PlayerView: MakePlayerView(roomSeat.player, false),
	}

	return roomSeatView
}

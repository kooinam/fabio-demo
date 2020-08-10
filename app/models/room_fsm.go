package models

import (
	"time"

	"github.com/kooinam/fabio/helpers"

	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio/models"
)

// RoomStates used to expose enums
var RoomStates *RoomFSMProperties

func init() {
	RoomStates = &RoomFSMProperties{
		Waiting:           "Waiting",
		Playing:           "Playing",
		Completed:         "Completed",
		PlayingDuration:   5 * time.Second,
		WaitingDuration:   0 * time.Second,
		CompletedDuration: 5 * time.Second,
	}
}

type RoomFSMProperties struct {
	Waiting           string
	Playing           string
	Completed         string
	PlayingDuration   time.Duration
	WaitingDuration   time.Duration
	CompletedDuration time.Duration
}

func makeRoomFSM(room *Room) *models.FiniteStateMachine {
	fsm := models.MakeFiniteStateMachine(RoomStates.Waiting)

	fsm.AddStateHandler(RoomStates.Waiting, room.enterWaiting, room.doWaiting, nil)
	fsm.AddStateHandler(RoomStates.Playing, room.enterPlaying, room.doPlaying, nil)
	fsm.AddStateHandler(RoomStates.Completed, room.enterCompleted, room.doCompleted, room.exitCompleted)

	return fsm
}

func (room *Room) enterWaiting(previous string) {
	roomView := MakeRoomView(room, true)

	fab.BroadcastEvent("room", room.ID, RoomStates.Waiting, roomView, nil)
}

func (room *Room) doWaiting() {
	emptySeat := room.Seats.Find(func(base models.Base) bool {
		seat := base.(*RoomSeat)

		return seat.isEmpty()
	})

	if emptySeat == nil {
		firstSeat := room.Seats.First()

		room.State.SetTurn(firstSeat, time.Now().Add(RoomStates.WaitingDuration))

		room.State.GoTo(RoomStates.Playing, room)
	}
}

func (room *Room) enterPlaying(previous string) {
	activeSeat := room.State.GetActiveAgent().(*RoomSeat)

	roomView := MakeRoomView(room, true)
	parameters := fab.H{
		"seat": activeSeat.Position,
	}

	fab.BroadcastEvent("room", room.ID, RoomStates.Playing, roomView, parameters)
}

func (room *Room) doPlaying() {
	activeSeat := room.State.GetActiveAgent().(*RoomSeat)

	if room.State.IsTurnExpired() {
		if activeSeat.hasMadeMove == false {
			// player has not make any move, make move for player

			for x := range room.Cells {
				for y := range room.Cells[x] {
					if room.Cells[x][y] == -1 {
						room.MakeMove(activeSeat.player, x, y)

						break
					}
				}

				if activeSeat.hasMadeMove {
					break
				}
			}
		}
	}

	if activeSeat.hasMadeMove {
		activeSeat.reset()

		winner := room.getWinner()

		if winner == -2 {
			// no winner, proceed to next seat
			nextSeat := room.getNextSeat()
			room.State.SetTurn(nextSeat, time.Now().Add(RoomStates.PlayingDuration))

			room.State.GoTo(RoomStates.Playing, room)
		} else {
			// has winner
			room.State.SetTurn(nil, time.Now().Add(RoomStates.PlayingDuration))

			room.State.GoTo(RoomStates.Completed, room)
		}
	}
}

func (room *Room) enterCompleted(previous string) {
	// increment winner in ranking
	winner := room.getWinner()
	winnerSeat, asserted := room.Seats.Find(func(base models.Base) bool {
		seat := base.(*RoomSeat)

		return seat.Position == winner
	}).(*RoomSeat)

	if asserted {
		room.Rankings[winnerSeat.player.Name]++
	}

	// broadcast completed event
	roomView := MakeRoomView(room, true)
	parameters := fab.H{
		"winner": winner,
	}

	fab.BroadcastEvent("room", room.ID, RoomStates.Completed, roomView, parameters)
}

func (room *Room) doCompleted() {
	if room.State.IsTurnExpired() {
		room.State.SetTurn(nil, time.Now().Add(RoomStates.CompletedDuration))

		room.State.GoTo(RoomStates.Waiting, room)
	}
}

func (room *Room) exitCompleted() {
	room.resetCells()
}

func (room *Room) getWinner() int {
	// -2 -> playing
	// -1 -> draw
	// 0 -> Player 1 wins
	// 1 -> Player 2 wins
	winner := -2

	helpers.Times(gridSize, func(x int) bool {
		if room.Cells[x][0] != -1 && room.Cells[x][0] == room.Cells[x][1] && room.Cells[x][1] == room.Cells[x][2] {
			winner = room.Cells[x][0]

			return false
		}

		return true
	})

	if winner == -2 {
		helpers.Times(gridSize, func(y int) bool {
			if room.Cells[0][y] != -1 && room.Cells[0][y] == room.Cells[1][y] && room.Cells[1][y] == room.Cells[2][y] {
				winner = room.Cells[0][y]

				return false
			}

			return true
		})
	}

	if winner == -2 {
		if room.Cells[0][0] != -1 && room.Cells[0][0] == room.Cells[1][1] && room.Cells[1][1] == room.Cells[2][2] {
			winner = room.Cells[0][0]
		}
	}

	if winner == -2 {
		if room.Cells[0][2] != -1 && room.Cells[0][2] == room.Cells[1][1] && room.Cells[1][1] == room.Cells[2][0] {
			winner = room.Cells[0][2]
		}
	}

	if winner == -2 {
		winner = -1

		for x := range room.Cells {
			for y := range room.Cells[x] {
				if room.Cells[x][y] == -1 {
					winner = -2

					break
				}
			}

			if winner == -2 {
				break
			}
		}
	}
	// if room.Cells[1][1] != -1 {
	// 	winner = room.Cells[1][1]
	// }

	return winner
}

func (room *Room) getNextSeat() *RoomSeat {
	activeSeat := room.State.GetActiveAgent().(*RoomSeat)

	nextSeat := room.Seats.Find(func(base models.Base) bool {
		seat := base.(*RoomSeat)

		return seat.Position > activeSeat.Position
	})

	if nextSeat == nil {
		// cant find seat's position that is greater active seat's position
		nextSeat = room.Seats.First().(*RoomSeat)
	}

	return nextSeat.(*RoomSeat)
}

package models

import (
	"github.com/kooinam/fabio/helpers"
)

// PlayerView used to represent player's view
type PlayerView struct {
	*Player
	AuthenticationToken string `json:"authenticationToken,omitempty"`
}

// MakePlayerView used to instantiate player's view
func MakePlayerView(player *Player, includeRoot bool) interface{} {
	var view interface{}

	if player != nil {
		// only marshal json if player is not nil, return nil otherwise
		playerView := &PlayerView{
			Player: player,
		}

		view = helpers.IncludeRootInJSON(playerView, includeRoot, "player")
	}

	return view
}

// MakeAuthenticatedPlayerView used to instantiate authenticated player's view
func MakeAuthenticatedPlayerView(player *Player, includeRoot bool) interface{} {
	var view interface{}

	if player != nil {
		// only marshal json if player is not nil, return nil otherwise
		playerView := &PlayerView{
			Player:              player,
			AuthenticationToken: player.AuthenticationToken,
		}

		view = helpers.IncludeRootInJSON(playerView, includeRoot, "player")
	}

	return view
}

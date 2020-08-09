package models

import (
	"github.com/kooinam/fabio/helpers"
)

func MakePlayerView(player *Player, includeRoot bool) interface{} {
	var view interface{}

	if player != nil {
		// only marshal json if player is not nil, return nil otherwise
		playerView := &struct {
			*Player
		}{
			Player: player,
		}

		view = helpers.IncludeRootInJSON(playerView, includeRoot, "player")
	}

	return view
}

func MakeAuthenticationPlayerView(player *Player, includeRoot bool) interface{} {
	var view interface{}

	if player != nil {
		// only marshal json if player is not nil, return nil otherwise
		playerView := &struct {
			*Player
			AuthenticationToken string `json:"authenticationToken"`
		}{
			Player:              player,
			AuthenticationToken: player.authenticationToken,
		}

		view = helpers.IncludeRootInJSON(playerView, includeRoot, "player")
	}

	return view
}

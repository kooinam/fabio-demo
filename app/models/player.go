package models

import (
	"github.com/kooinam/fabio/models"
	"syreclabs.com/go/faker"
)

// PlayersCollection is singleton for RoomsCollection
var PlayersCollection *models.Collection

// Player used to represents player
type Player struct {
	models.Base
	Name                string
	authenticationToken string
}

// MakePlayer used to instantiate player
func MakePlayer(collection *models.Collection, args ...interface{}) models.Modellable {
	player := &Player{}

	player.Initialize(collection)

	player.Name = faker.Internet().UserName()
	player.authenticationToken = faker.Internet().Password(24, 24)

	return player
}

// GetAuthenticationToken used to get authentication token
func (player *Player) GetAuthenticationToken() string {
	return player.authenticationToken
}

// AuthenticatePlayer used to authenticate a player
func AuthenticatePlayer(token string) *Player {
	authenticatedPlayer := PlayersCollection.FindOrCreate(func(item models.Modellable) bool {
		player := item.(*Player)

		return player.authenticationToken == token
	}).(*Player)

	return authenticatedPlayer
}

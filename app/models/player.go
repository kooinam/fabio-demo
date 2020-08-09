package models

import (
	"fmt"

	"github.com/kooinam/fabio/models"
	"syreclabs.com/go/faker"
)

// PlayersCollection is singleton for RoomsCollection
var PlayersCollection *models.Collection

func init() {
	PlayersCollection = models.MakeCollection(makePlayer)
}

// Player used to represents player
type Player struct {
	ID                  string `json:"id"`
	Name                string
	authenticationToken string
}

func makePlayer(args ...interface{}) models.Base {
	player := &Player{
		ID:                  fmt.Sprintf("%v", PlayersCollection.Count()+1),
		Name:                faker.Internet().UserName(),
		authenticationToken: faker.Internet().Password(24, 24),
	}

	return player
}

// GetID used to get ID
func (player *Player) GetID() string {
	return player.ID
}

// GetAuthenticationToken used to get authentication token
func (player *Player) GetAuthenticationToken() string {
	return player.authenticationToken
}

// AuthenticatePlayer used to authenticate a player
func AuthenticatePlayer(token string) *Player {
	authenticatedPlayer := PlayersCollection.FindOrCreate(func(base models.Base) bool {
		player := base.(*Player)

		return player.authenticationToken == token
	}).(*Player)

	return authenticatedPlayer
}

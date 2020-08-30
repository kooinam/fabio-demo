package models

import (
	"github.com/kooinam/fabio/helpers"
	"github.com/kooinam/fabio/models"
	"github.com/kooinam/fabio/mongorecords"
	"syreclabs.com/go/faker"
)

// PlayersCollection is singleton for RoomsCollection
var PlayersCollection *models.Collection

// Player used to represents player
type Player struct {
	mongorecords.Base   `bson:"base,inline"`
	Name                string `json:"name" bson:"name"`
	AuthenticationToken string `json:"-" bson:"authentication_token"`
}

// FindPlayerByToken used to find a player by token
func FindPlayerByToken(token string) *models.SingleResult {
	result := PlayersCollection.Query().Where(helpers.H{
		"authentication_token": token,
	}).First()

	return result
}

// MakePlayer used to instantiate player
func MakePlayer(collection *models.Collection, hooksHandler *models.HooksHandler) models.Modellable {
	player := &Player{}

	hooksHandler.RegisterInitializeHook(player.initialize)

	return player
}

func (player *Player) initialize(attributes *helpers.Dictionary) {
	player.Name = faker.Internet().UserName()
	player.AuthenticationToken = faker.Internet().Password(24, 24)
}

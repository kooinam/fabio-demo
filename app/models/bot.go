package models

import (
	fab "github.com/kooinam/fabio"
	"github.com/kooinam/fabio/actors"
	"github.com/kooinam/fabio/helpers"
	"github.com/kooinam/fabio/models"
	"github.com/kooinam/fabio/simplerecords"
)

// BotsCollection used to retrieve registered bots collections
func BotsCollection() *models.Collection {
	return fab.ModelManager().Collection("simple", "bots")
}

type Bot struct {
	simplerecords.Base
	actor *actors.Actor
	timer float64
}

func MakeBot(collection *models.Collection, hooksHandler *models.HooksHandler) models.Modellable {
	bot := &Bot{}

	hooksHandler.RegisterInitializeHook(bot.initialize)

	hooksHandler.RegisterAfterMemoizeHook(bot.afterMemoize)

	return bot
}

func (bot *Bot) initialize(attributes *helpers.Dictionary) {}

func (bot *Bot) afterMemoize() {
	bot.actor = fab.ActorManager().RegisterActor(bot.GetCollectionName(), bot)
}

// RegisterActorActions used to register actor's actions
func (bot *Bot) RegisterActorActions(actionsHandler *actors.ActionsHandler) {
	actionsHandler.RegisterAction("Update", bot.update)
}

func (bot *Bot) update(context *actors.Context) error {
	var err error

	dt := context.ParamsFloat64("dt", 0)
	_ = dt

	return err
}

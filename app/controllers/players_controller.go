package controllers

import (
	"github.com/kooinam/fabio-demo/app/models"
	"github.com/kooinam/fabio/controllers"
	"github.com/kooinam/fabio/helpers"
)

// PlayersController is controller for player's actions
type PlayersController struct {
}

// RegisterHooksAndActions used to register hooks and actions
func (controller *PlayersController) RegisterHooksAndActions(hooksHandler *controllers.HooksHandler, actionsHandler *controllers.ActionsHandler) {
	actionsHandler.RegisterAction("Register", controller.register)
}

func (controller *PlayersController) register(context *controllers.Context) {
	result := models.PlayersCollection.Create(helpers.H{})

	if result.StatusError() {
		context.SetErrorResult(controllers.StatusError, result.Error())
		return
	}

	player := result.Item().(*models.Player)
	playerView := models.MakeAuthenticatedPlayerView(player, true)

	context.SetSuccessResult(playerView)
	return
}

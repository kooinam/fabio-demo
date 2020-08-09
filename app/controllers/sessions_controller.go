package controllers

import (
	"github.com/kooinam/fabio-demo/app/models"

	"github.com/kooinam/fabio/controllers"
)

// SessionsController used for sessions actions
type SessionsController struct {
}

// AddBeforeActions used to add before actions callbacks
func (controller *SessionsController) AddBeforeActions(callbacksHandler *controllers.CallbacksHandler) {
	// controller.callbacksHandler.AddBeforeAction(controller.SetCurrentPlayer)
}

// AddActions used to add actions
func (controller *SessionsController) AddActions(actionsHandler *controllers.ActionsHandler) {
	actionsHandler.AddAction("Authenticate", controller.authenticate)
}

// authenticate used to authenticate a player
func (controller *SessionsController) authenticate(connection *controllers.Connection) (interface{}, error) {
	var playerView interface{}

	authenticationToken := connection.ParamsWithFallback("authenticationToken", "").(string)

	player := models.AuthenticatePlayer(authenticationToken)
	playerView = models.MakeAuthenticationPlayerView(player, true)

	return playerView, nil
}

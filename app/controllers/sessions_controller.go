package controllers

import (
	"github.com/kooinam/fabio-demo/app/models"

	"github.com/kooinam/fabio/controllers"
)

// SessionsController used for sessions actions
type SessionsController struct {
}

// RegisterBeforeHooks used to register before action hooks
func (controller *SessionsController) RegisterBeforeHooks(hooksHandler *controllers.HooksHandler) {
}

// RegisterActions used to register actions
func (controller *SessionsController) RegisterActions(actionsHandler *controllers.ActionsHandler) {
	actionsHandler.RegisterAction("Authenticate", controller.authenticate)
}

// authenticate used to authenticate a player
func (controller *SessionsController) authenticate(connection *controllers.Context) (interface{}, error) {
	var playerView interface{}

	authenticationToken := connection.ParamsStr("authenticationToken")

	player := models.AuthenticatePlayer(authenticationToken)
	playerView = models.MakeAuthenticatedPlayerView(player, true)

	return playerView, nil
}

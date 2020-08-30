package controllers

import (
	"fmt"

	"github.com/kooinam/fabio-demo/app/models"

	"github.com/kooinam/fabio/controllers"
)

// SessionsController used for sessions actions
type SessionsController struct {
}

// RegisterHooksAndActions used to register hooks and actions
func (controller *SessionsController) RegisterHooksAndActions(hooksHandler *controllers.HooksHandler, actionsHandler *controllers.ActionsHandler) {
	actionsHandler.RegisterAction("Authenticate", controller.authenticate)
}

// authenticate used to authenticate a player
func (controller *SessionsController) authenticate(context *controllers.Context) {
	authenticationToken := context.ParamsStr("authenticationToken")

	result := models.FindPlayerByToken(authenticationToken)

	if !result.StatusSuccess() {
		err := fmt.Errorf("unauthorized")
		context.SetErrorResult(controllers.StatusUnauthorized, err)
		return
	}

	playerView := models.MakeAuthenticatedPlayerView(result.Item().(*models.Player), true)

	context.SetSuccessResult(playerView)
}

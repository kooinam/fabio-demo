package controllers

import "github.com/kooinam/fabio/controllers"

type PlayersController struct {
}

// AddBeforeActions used to add before actions callbacks
func (controller *PlayersController) AddBeforeActions(callbacksHandler *controllers.CallbacksHandler) {
	// controller.callbacksHandler.AddBeforeAction(controller.SetCurrentPlayer)
}

// AddActions used to add actions
func (controller *PlayersController) AddActions(actionsHandler *controllers.ActionsHandler) {
	// controller.callbacksHandler.AddBeforeAction(controller.SetCurrentPlayer)
}

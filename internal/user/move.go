package user

import (
	"net/http"
	"startrader/internal/response"
	"startrader/internal/starsystem"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MoveInput struct {
	System string `form:"system"`
}

func MovePostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("id")

	var moveInput MoveInput
	if c.ShouldBind(&moveInput) != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Cannot bind form data",
		}, http.StatusBadRequest)
		return
	}

	users, err := ReadUsers([]string{userID.(string)})
	u := users[userID.(string)]

	if err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Can't find user",
			Data:    userID,
		}, http.StatusInternalServerError)
		return
	}

	destination, err := starsystem.ReadSystem(moveInput.System)
	if err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Unknown destination",
			Data: map[string]any{
				"destination": moveInput.System,
			},
		}, http.StatusBadRequest)
		return
	}

	fuelRequired, _ := u.FuelRequired(destination)

	fuelAvailable := float64(0)
	fuel, ok := u.Inventory["fuel"]
	if ok {
		fuelAvailable = fuel.Quantity
		fuel.Quantity -= fuelRequired
	}
	if fuelAvailable < fuelRequired {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Not enough fuel to move",
			Data: map[string]any{
				"required_fuel":  fuelRequired,
				"available_fuel": fuelAvailable,
				"destination":    destination.Name,
			},
		}, http.StatusOK)
		return
	}

	// Deregister trader from the source system
	if err := starsystem.DeregisterTrader(u.Name, u.Location); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Failed to deregister user from the starbase",
			Data:    map[string]any{},
		}, http.StatusInternalServerError)
		return
	}

	// Update user data
	u.Location = destination.Name
	u.Inventory["fuel"] = fuel
	if err := WriteUserState(&u); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Failed to save state",
			Data:    map[string]any{},
		}, http.StatusInternalServerError)
		return
	}

	// Register trader in the destination system
	if err := starsystem.RegisterTrader(u.Name, u.Location); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Failed to register user in the starbase",
			Data:    map[string]any{},
		}, http.StatusInternalServerError)
		return
	}

	// Return successful status to the caller
	d := map[string]interface{}{
		"user":           u.Name,
		"start":          u.Location,
		"destination":    destination.Name,
		"required_fuel":  fuelRequired,
		"remaining_fuel": u.Inventory["fuel"].Quantity,
	}
	response.WriteResponse(c, response.Response{
		Status:  response.Ok,
		Message: "Moved",
		Data:    d,
	}, http.StatusCreated)
}

package user

import (
	"net/http"
	"startrader/internal/auth"
	"startrader/internal/response"
	"startrader/internal/starsystem"
)

func MovePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.SessionData(r, "userID")
	u, _ := ReadUser(userID.(string))
	destination, _ := starsystem.ReadSystem(r.FormValue("system"))

	fuelRequired, _ := u.FuelRequired(destination)

	fuelAvailable := float64(0)
	fuel, ok := u.Inventory["fuel"]
	if ok {
		fuelAvailable = fuel.Quantity
		fuel.Quantity -= fuelRequired
	}
	if fuelAvailable < fuelRequired {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: "Not enough fuel to move",
			Data: map[string]any{
				"required_fuel":  fuelRequired,
				"available_fuel": fuelAvailable,
				"destination":    destination.Name,
			},
		})
		return
	}

	u.Location = destination.Name
	u.Inventory["fuel"] = fuel
	err := WriteUserState(&u)
	if err != nil {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: "Failed to save state",
			Data:    map[string]any{},
		})
		return
	}

	d := map[string]interface{}{
		"user":           u.Name,
		"start":          u.Location,
		"destination":    destination.Name,
		"required_fuel":  fuelRequired,
		"remaining_fuel": u.Inventory["fuel"].Quantity,
	}
	response.WriteResponse(w, response.Response{
		Status:  response.Ok,
		Message: "Moved",
		Data:    d,
	})
}

package user

import (
	"fmt"
	"log"
	"net/http"
	"startrader/internal/auth"
	"startrader/internal/response"
	"startrader/internal/starsystem"
	"strconv"
)

func BuyPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("BuyHandler")
	userID := auth.SessionData(r, "userID")
	u, _ := ReadUser(userID.(string))
	s, err := starsystem.ReadSystem(u.Location)
	if err != nil {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: err.Error(),
		})
		return
	}

	item := r.FormValue("item")
	quantity, err := strconv.ParseFloat(r.FormValue("quantity"), 64)
	if err != nil {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: err.Error(),
			Data: map[string]any{
				"input_value": r.FormValue("quantity"),
			},
		})
		return
	}

	marketItem, ok := s.Market[item]
	if !ok {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: "Item not found",
			Data: map[string]any{
				"item": item,
			},
		})
		return
	}

	if marketItem.Quantity < quantity {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: "Not enough available items",
			Data: map[string]any{
				"item":               item,
				"requested_quantity": fmt.Sprintf("%.1f", quantity),
				"available_quantity": fmt.Sprintf("%.1f", marketItem.Quantity),
			},
		})
		return
	}

	requiredCredits := marketItem.Price * quantity
	if u.Credits < requiredCredits {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: "Not enough credits",
			Data: map[string]any{
				"required_credits":  fmt.Sprintf("%.1f", requiredCredits),
				"available_credits": fmt.Sprintf("%.1f", u.Credits),
			},
		})
		return
	}

	u.Credits -= requiredCredits
	u.Inventory.Add(item, quantity)
	s.Market.Reduce(item, quantity)

	if err := starsystem.WriteSystemState(s); err != nil {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: err.Error(),
		})
		return
	}

	if err := WriteUserState(&u); err != nil {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: err.Error(),
		})
		return
	}

	response.WriteResponse(w, response.Response{
		Status:  response.Ok,
		Message: "Bought",
		Data: map[string]any{
			"item":     item,
			"quantity": fmt.Sprintf("%.1f", quantity),
		},
	})
}

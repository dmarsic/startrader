package user

import (
	"fmt"
	"log"
	"net/http"
	"startrader/internal/response"
	"startrader/internal/starsystem"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type BuyInput struct {
	Item     string  `form:"item"`
	Quantity float64 `form:"quantity"`
}

func BuyPostHandler(c *gin.Context) {
	log.Println("BuyHandler")

	session := sessions.Default(c)
	userID := session.Get("id")

	users, err := ReadUsers([]string{userID.(string)})
	if err != nil || len(users) == 0 {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Can't read user data",
			Data: map[string]any{
				"user": userID,
			},
		}, http.StatusInternalServerError)
	}
	u := users[userID.(string)]

	s, err := starsystem.ReadSystem(u.Location)
	if err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	var buyInput BuyInput
	if c.ShouldBind(&buyInput) != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Cannot bind form data",
		}, http.StatusBadRequest)
		return
	}

	marketItem, ok := s.Market[buyInput.Item]
	if !ok {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Item not found",
			Data: map[string]any{
				"item": buyInput.Item,
			},
		}, http.StatusOK)
		return
	}

	if marketItem.Quantity < buyInput.Quantity {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Not enough available items",
			Data: map[string]any{
				"item":               buyInput.Item,
				"requested_quantity": fmt.Sprintf("%.1f", buyInput.Quantity),
				"available_quantity": fmt.Sprintf("%.1f", marketItem.Quantity),
			},
		}, http.StatusOK)
		return
	}

	requiredCredits := marketItem.Price * buyInput.Quantity
	if u.Credits < requiredCredits {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Not enough credits",
			Data: map[string]any{
				"required_credits":  fmt.Sprintf("%.1f", requiredCredits),
				"available_credits": fmt.Sprintf("%.1f", u.Credits),
			},
		}, http.StatusOK)
		return
	}

	u.Credits -= requiredCredits
	u.Inventory.Add(buyInput.Item, buyInput.Quantity)
	s.Market.Reduce(buyInput.Item, buyInput.Quantity)

	if err := starsystem.WriteSystemState(s); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := WriteUserState(&u); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	response.WriteResponse(c, response.Response{
		Status:  response.Ok,
		Message: "Bought",
		Data: map[string]any{
			"item":     buyInput.Item,
			"quantity": fmt.Sprintf("%.1f", buyInput.Quantity),
		},
	}, http.StatusCreated)
}

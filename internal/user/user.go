package user

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"startrader/internal/response"
	"startrader/internal/starsystem"

	"github.com/abhinav-TB/dantdb"
	"github.com/gin-gonic/gin"
)

const StoreDir = "./data/"
const StartingCredits = 1000.0
const StartingSystem = "sol"

type UserList struct {
	Users map[string]User
}

type Users map[string]User

type User struct {
	Name      string    `json:"name"`
	Credits   float64   `json:"credits"`
	Location  string    `json:"location"`
	Inventory Inventory `json:"inventory"`
}

type Inventory map[string]Item

type Item struct {
	Quantity float64 `json:"quantity"`
}

type UserInput struct {
	Name string `json:"name"`
}

func NewUser(name string) *User {
	return &User{
		Name:      name,
		Credits:   StartingCredits,
		Location:  StartingSystem,
		Inventory: nil,
	}
}

func NewUserPostHandler(c *gin.Context) {
	// Struct to bind the request body
	var userInput UserInput

	// Bind the request body
	if err := c.ShouldBindJSON(&userInput); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: err.Error(),
			Data:    map[string]any{},
		}, http.StatusBadRequest)
		return
	}

	name := userInput.Name
	log.Println("user.NewUserPostHandler: name=" + name)

	// If user exists, return error.
	_, err := ReadUsers([]string{name})
	if err == nil {
		// User exists.
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "User exists",
			Data: map[string]any{
				"name": name,
			},
		}, http.StatusOK)
		return
	}

	// Create new user.
	u := NewUser(name)
	err = WriteUserState(u)
	if err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Failed to save state",
			Data:    map[string]any{},
		}, http.StatusInternalServerError)
		return
	}

	// Register as trader in a starbase.
	if err := starsystem.RegisterTrader(name, u.Location); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "Failed to save state",
			Data:    map[string]any{},
		}, http.StatusInternalServerError)
	}

	// Success response.
	response.WriteResponse(c, response.Response{
		Status:  response.Ok,
		Message: "Created",
		Data: map[string]any{
			"name": name,
		},
	}, http.StatusCreated)

	log.Println("user.NewUserPostHandler: created new user: " + name)
}

func (u User) FuelRequired(destination *starsystem.System) (float64, error) {
	currentSystem, err := starsystem.ReadSystem(u.Location)
	if err != nil {
		return 0.0, err
	}
	deltaX := float64(destination.Position.X) - float64(currentSystem.Position.X)
	deltaY := float64(destination.Position.Y) - float64(currentSystem.Position.Y)
	deltaZ := float64(destination.Position.Z) - float64(currentSystem.Position.Z)
	return math.Sqrt(math.Pow(deltaX, 2.0) + math.Pow(deltaY, 2.0) + math.Pow(deltaZ, 2.0)), nil
}

func ReadAllUsers() (Users, error) {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return nil, err
	}

	records, err := db.ReadAll("users")
	if err != nil {
		return nil, err
	}

	users := Users{}
	for _, r := range records {
		user := User{}
		if err := json.Unmarshal([]byte(r), &user); err != nil {
			return nil, err
		}
		users[user.Name] = user
	}

	return users, nil
}

func ReadUsers(names []string) (Users, error) {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return nil, err
	}

	records, err := db.ReadFiltered("users", names)
	if err != nil {
		return nil, err
	}

	users := Users{}
	for _, r := range records {
		user := User{}
		if err := json.Unmarshal([]byte(r), &user); err != nil {
			return nil, err
		}
		users[user.Name] = user
	}

	return users, nil
}

func WriteUserState(u *User) error {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return err
	}

	if err := db.Write("users", u.Name, u); err != nil {
		return err
	}

	return nil
}

// Add adds an item to inventory.
func (i *Inventory) Add(item string, quantity float64) {
	existing, ok := (*i)[item]
	if !ok {
		(*i)[item] = Item{
			Quantity: quantity,
		}
	} else {
		newQuantity := existing.Quantity + quantity
		(*i)[item] = Item{
			Quantity: newQuantity,
		}
	}
}

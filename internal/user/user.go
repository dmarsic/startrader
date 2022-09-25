package user

import (
	"encoding/json"
	"math"
	"startrader/internal/starsystem"

	"github.com/abhinav-TB/dantdb"
)

const StoreDir = "./data/"

type UserList struct {
	Users map[string]User
}

type Users map[string]User

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Credits   float64   `json:"credits"`
	Location  string    `json:"location"`
	Inventory Inventory `json:"inventory"`
}

type Inventory map[string]Item

type Item struct {
	Quantity float64 `json:"quantity"`
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

func ReadUser(userID string) (User, error) {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = db.Read("users", userID, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
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

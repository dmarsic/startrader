package main

import (
	"fmt"
	"log"
	"net/http"

	"startrader/internal/auth"
	"startrader/internal/response"
	"startrader/internal/starsystem"
	"startrader/internal/user"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/systems", AllSystemsHandler)
	r.HandleFunc("/systems/{name}", SystemGetHandler)
	r.HandleFunc("/u", AllUsersHandler)
	r.HandleFunc("/u/new", user.NewUserPostHandler).Methods(http.MethodPost)
	r.HandleFunc("/u/{name}", UserGetHandler)
	r.HandleFunc("/m", user.MovePostHandler).Methods(http.MethodPost)
	r.HandleFunc("/b", user.BuyPostHandler).Methods(http.MethodPost)

	r.Use(auth.AuthMiddleware)

	port := ":5000"
	fmt.Println("Server is running on port" + port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HomeHandler")
	userID := auth.SessionData(r, "userID")
	http.Redirect(w, r, fmt.Sprintf("/u/%s", userID), http.StatusFound)
}

func AllSystemsHandler(w http.ResponseWriter, r *http.Request) {
	systems, _ := starsystem.ReadAllSystems()
	response.WriteResponse(w, response.Response{
		Status: response.Ok,
		Data:   systems,
	})
}

func SystemGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, _ := starsystem.ReadSystem(vars["name"])
	response.WriteResponse(w, response.Response{
		Status: response.Ok,
		Data:   s,
	})
}

func AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AllUsersHandler")
	users, _ := user.ReadAllUsers()
	response.WriteResponse(w, response.Response{
		Status: response.Ok,
		Data:   users,
	})
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UserGetHandler")
	vars := mux.Vars(r)
	u, _ := user.ReadUser(vars["name"])
	response.WriteResponse(w, response.Response{
		Status: response.Ok,
		Data:   u,
	})
}

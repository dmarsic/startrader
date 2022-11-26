package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"startrader/internal/auth"
	"startrader/internal/response"
	"startrader/internal/starsystem"
	"startrader/internal/user"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	a := r.PathPrefix("/api/v1").Subrouter()
	a.HandleFunc("/", HomeHandler)
	a.HandleFunc("/login", auth.LoginHandler)
	a.HandleFunc("/logout", auth.LogoutHandler)
	a.HandleFunc("/systems", AllSystemsHandler)
	a.HandleFunc("/systems/{name}", SystemGetHandler)
	a.HandleFunc("/u", AllUsersHandler)
	a.HandleFunc("/u/new", user.NewUserPostHandler).Methods(http.MethodPost)
	a.HandleFunc("/u/{name}", UserGetHandler)
	a.HandleFunc("/m", user.MovePostHandler).Methods(http.MethodPost)
	a.HandleFunc("/b", user.BuyPostHandler).Methods(http.MethodPost)

	a.Use(auth.AuthMiddleware)

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
	}, http.StatusOK)
}

func SystemGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, _ := starsystem.ReadSystem(vars["name"])
	response.WriteResponse(w, response.Response{
		Status: response.Ok,
		Data:   s,
	}, http.StatusOK)
}

func AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AllUsersHandler")
	users, _ := user.ReadAllUsers()
	response.WriteResponse(w, response.Response{
		Status: response.Ok,
		Data:   users,
	}, http.StatusOK)
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UserGetHandler")

	vars := mux.Vars(r)
	userList := strings.Split(vars["name"], ",")
	if len(userList) == 0 || userList[0] == "" {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: "User(s) not specified",
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	u, err := user.ReadUsers(userList)
	if err != nil {
		response.WriteResponse(w, response.Response{
			Status:  response.Error,
			Message: err.Error(),
			Data:    userList,
		}, http.StatusInternalServerError)
	}
	response.WriteResponse(w, response.Response{
		Status: response.Ok,
		Data:   u,
	}, http.StatusOK)
}

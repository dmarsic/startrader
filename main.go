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

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("hushhush"))
	r.Use(sessions.Sessions("session", store))

	a := r.Group("/api/v1")
	{
		a.GET("/", HomeHandler)
		a.GET("/login", auth.LoginHandler)
		a.GET("/logout", auth.LogoutHandler)
		a.GET("/systems", AllSystemsHandler)
		a.GET("/systems/:name", SystemGetHandler)
		a.GET("/u", AllUsersHandler)
		a.POST("/u/new", user.NewUserPostHandler)
		a.GET("/u/:name", UserGetHandler)
		a.POST("/m", user.MovePostHandler)
		a.POST("/b", user.BuyPostHandler)
	}

	a.Use(auth.AuthMiddleware())

	port := ":9500"
	fmt.Println("Server is running on port" + port)

	r.Run(port)
}

func HomeHandler(c *gin.Context) {
	log.Println("HomeHandler")
	session := sessions.Default(c)
	userID := session.Get("id")
	group := c.FullPath()[:len(c.FullPath())-1]
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/u/%s", group, userID))
}

func AllSystemsHandler(c *gin.Context) {
	systems, _ := starsystem.ReadAllSystems()
	response.WriteResponse(c, response.Response{
		Status: response.Ok,
		Data:   systems,
	}, http.StatusOK)
}

func SystemGetHandler(c *gin.Context) {
	name := c.Param("name")
	s, _ := starsystem.ReadSystem(name)
	response.WriteResponse(c, response.Response{
		Status: response.Ok,
		Data:   s,
	}, http.StatusOK)
}

func AllUsersHandler(c *gin.Context) {
	log.Println("AllUsersHandler")
	users, _ := user.ReadAllUsers()
	response.WriteResponse(c, response.Response{
		Status: response.Ok,
		Data:   users,
	}, http.StatusOK)
}

type UserInput struct {
	Name string `uri:"name"`
}

func UserGetHandler(c *gin.Context) {
	log.Println("UserGetHandler")

	var userInput UserInput

	if err := c.ShouldBindUri(&userInput); err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: err.Error(),
			Data:    map[string]any{},
		}, http.StatusBadRequest)
		return
	}

	name := userInput.Name

	userList := strings.Split(name, ",")
	if len(userList) == 0 || userList[0] == "" {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: "User(s) not specified",
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	u, err := user.ReadUsers(userList)
	if err != nil {
		response.WriteResponse(c, response.Response{
			Status:  response.Error,
			Message: err.Error(),
			Data:    userList,
		}, http.StatusInternalServerError)
	}
	response.WriteResponse(c, response.Response{
		Status: response.Ok,
		Data:   u,
	}, http.StatusOK)
}

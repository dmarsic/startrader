// // Package auth manages user login and logout.
// // It stores userID into a cookie to use it as session data.
// //
// // There is no real authentication. The user can just state
// // who they are at this time, and it will be stored in the
// // cookie.

package auth

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"startrader/internal/response"
)

func LoginHandler(c *gin.Context) {
	userID := c.Query("user")
	log.Printf("LoginHandler: user=%s\n", userID)

	session := sessions.Default(c)
	session.Set("id", userID)
	session.Save()
	response.WriteResponse(c, response.Response{
		Status:  response.Ok,
		Message: "Logged in",
		Data: map[string]any{
			"user": userID,
		},
	}, http.StatusOK)
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("id")
	session.Clear()
	session.Save()

	log.Printf("LogoutHandler: user=%s\n", sessionID)
	response.WriteResponse(c, response.Response{
		Status:  response.Ok,
		Message: "Logged out",
		Data: map[string]any{
			"userID": sessionID,
		},
	}, http.StatusOK)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if urlPath == "/api/v1/login" || urlPath == "/api/v1/logout" || urlPath == "/api/v1/u/new" {
			log.Println("authMiddleware: matched path: " + urlPath + ", skipping login redirect.")
			c.Next()
			return
		}

		session := sessions.Default(c)
		sessionID := session.Get("id")
		if sessionID == nil {
			response.WriteResponse(c, response.Response{
				Status:  response.Error,
				Message: "Not logged in",
				Data:    map[string]any{},
			}, http.StatusUnauthorized)
			c.Redirect(http.StatusFound, "/api/v1/login")
			return
		}
		log.Printf("authMiddleware: session id=%s, passing\n", sessionID)
	}
}

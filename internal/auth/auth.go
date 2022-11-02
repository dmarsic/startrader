// Package auth manages user login and logout.
// It stores userID into a cookie to use it as session data.
//
// There is no real authentication. The user can just state
// who they are at this time, and it will be stored in the
// cookie.

package auth

import (
	"log"
	"net/http"
	"os"
	"startrader/internal/response"

	"github.com/gorilla/sessions"
)

var cookiejar = sessions.NewCookieStore([]byte(os.Getenv("STARTRADER_SESSION_KEY")))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	log.Printf("LoginHandler: user=%s\n", userID)

	session, _ := cookiejar.Get(r, "session")
	session.Values["userID"] = userID
	session.Save(r, w)

	response.WriteResponse(w, response.Response{
		Status:  response.Ok,
		Message: "Logged in",
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := cookiejar.Get(r, "session")
	userID, ok := session.Values["userID"]
	log.Printf("LogoutHandler: user=%s\n", userID)
	if !ok {
		response.WriteResponse(w, response.Response{
			Status:  response.Warning,
			Message: "No user logged in",
		})
	} else {
		delete(session.Values, "userID")
		session.Save(r, w)
		response.WriteResponse(w, response.Response{
			Status:  response.Ok,
			Message: "Logged out",
			Data: map[string]any{
				"userID": userID,
			},
		})
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/logout" || r.URL.Path == "/u/new" {
			log.Println("authMiddleware: matched path: " + r.URL.Path + ", skipping login redirect.")
			next.ServeHTTP(w, r)
			return
		}

		userID := SessionData(r, "userID")
		if userID == nil {
			log.Println("authMiddleware: userID is not set, going to /login")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		log.Printf("authMiddleware: user=%s, passing\n", userID)
		next.ServeHTTP(w, r)
	})
}

func SessionData(r *http.Request, key string) interface{} {
	session, _ := cookiejar.Get(r, "session")
	value, ok := session.Values[key]
	if !ok {
		return nil
	}
	return value
}

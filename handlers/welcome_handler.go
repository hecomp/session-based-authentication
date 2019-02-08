package handlers

import (
	"fmt"
	"github.com/hecomp/session-based-authentication/redisclient"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	response, err := redisclient.Cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		// I fthe session token is not in cache, return an unauthorized errror
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
	}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code up until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	response, err := redisclient.Cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code up until this point is the same as the first part of the `Welcome` route

	// Now, create a new session token for the the current user
	newSessionToken, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = redisclient.Cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s", response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the older session token
	_, err = redisclient.Cache.Do("DEL", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:"session_token",
		Value:newSessionToken.String(),
		Expires:time.Now().Add(120 * time.Second),
	})
}

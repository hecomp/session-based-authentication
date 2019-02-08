package handlers

import (
	"encoding/json"
	"github.com/hecomp/session-based-authentication/redisclient"
	"github.com/satori/go.uuid"
	"time"

	//"github.com/satori/go.uuid"
	"net/http"
)

var users = map[string]string{
	"user1":"password1",
	"user2":"password2",
	"user3":"password3",
}

type Credentials struct {
	Password string `json:"password""`
	Username string `json:"username"`
}

func Sigin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		//	If the structure of the body is wrong , return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}


	// Create a new random session token
	sessionToken, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Set the in the cache, along with the user whom it represents
	// The token has expiry time of 120
	_, err = redisclient.Cache.Do("SETEX", sessionToken, "120", creds.Username)
	if err != nil {
	//	If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for the "session_token" as the session token we just generated
	// e also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: sessionToken.String(),
		Expires:time.Now().Add(120 * time.Second),
	})


}
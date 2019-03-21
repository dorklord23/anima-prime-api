package middlewares

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dorklord23/anima-prime/routes"
	"github.com/dorklord23/anima-prime/utils"
)

// Authenticate : function to authenticate a request based on a particular header in the request
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		anonymousEndpoints := make(map[string][]string)
		anonymousEndpoints["POST"] = []string{"/api/users"}

		// Check if this request doesn't need to be authenticated based on the method...
		if len(anonymousEndpoints[r.Method]) > 0 {
			// ... and  the URL
			if utils.Contains(anonymousEndpoints[r.Method], r.RequestURI) {
				// No need for authentication
				next.ServeHTTP(w, r)
				return
			}
		}

		// Check if the request header exists
		if r.Header.Get("anima-prime-token") == "" {
			data := make(map[string]string)
			data["Message"] = "This request cannot be authenticated"
			routes.SendResponse(w, 401, data, "fail")
		} else {
			// Decode the token
			decodedTokenInBytes, err := base64.StdEncoding.DecodeString(r.Header.Get("anima-prime-token"))
			if err != nil {
				routes.SendResponse(w, 500, err.Error(), "error")
				return
			}

			// Check if it's already expired
			decodedToken := utils.BytesToString(decodedTokenInBytes)
			splitStrings := strings.Split(decodedToken, "|")
			creationDate := splitStrings[1]
			tokenAge, err3 := strconv.Atoi(splitStrings[2])
			if err3 != nil {
				routes.SendResponse(w, 500, err3.Error(), "error")
				return
			}

			t, err2 := time.Parse("2006-01-02T15:04:05.000Z", creationDate)
			if err2 != nil {
				routes.SendResponse(w, 500, err2.Error(), "error")
				return
			}

			expiryTime := t.AddDate(0, 0, tokenAge)

			if time.Now().After(expiryTime) {
				data := make(map[string]string)
				data["Message"] = "Your token has expired."
				routes.SendResponse(w, 401, data, "fail")
				return
			}

			next.ServeHTTP(w, r)
		}
	})
}

/* Copyright 2019 Tri Rumekso Anggie Wibowo (trirawibowo [at] gmail [dot] com)
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package middlewares

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dorklord23/anima-prime/models"
	"github.com/dorklord23/anima-prime/utils"
	"github.com/gorilla/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Authenticate : function to authenticate a request based on a particular header in the request
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		anonymousEndpoints := make(map[string][]string)
		anonymousEndpoints["POST"] = []string{
			"/api/users", "/api/login",
		}

		anonymousEndpoints["GET"] = []string{
			"/api/tokens",
		}

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
			utils.SendResponse(w, 401, data, "fail", nil)
		} else {
			// Decode the token
			decodedTokenInBytes, err := base64.StdEncoding.DecodeString(r.Header.Get("anima-prime-token"))
			if err != nil {
				data := make(map[string]string)
				data["Message"] = "Invalid access token"
				utils.SendResponse(w, 400, data, "fail", nil)
				return
			}

			// Check if it's already expired
			decodedToken := utils.BytesToString(decodedTokenInBytes)
			splitStrings := strings.Split(decodedToken, "|")
			creationDate := splitStrings[1]
			tokenAge, err3 := strconv.Atoi(splitStrings[2])
			if err3 != nil {
				utils.SendResponse(w, 500, err3.Error(), "error", nil)
				return
			}

			t, err2 := time.Parse("2006-01-02T15:04:05Z", creationDate)
			if err2 != nil {
				utils.SendResponse(w, 500, err2.Error(), "error", nil)
				return
			}

			expiryTime := t.AddDate(0, 0, tokenAge)

			if time.Now().After(expiryTime) {
				data := make(map[string]string)
				data["Message"] = "Your token has expired."
				utils.SendResponse(w, 401, data, "fail", nil)
				return
			}

			// Decode the key
			key, err4 := datastore.DecodeKey(splitStrings[0])
			if err4 != nil {
				utils.SendResponse(w, 500, err4.Error(), "error", nil)
				return
			}

			// Pass the requester's user key and authority to the handler
			var userStruct models.User
			err5 := datastore.Get(ctx, key, &userStruct)
			if err5 == datastore.Done {
				// The requester is not a registered user
				data := make(map[string]string)
				data["Email"] = "The requester is not recognized"
				utils.SendResponse(w, 401, data, "fail", nil)
				return
			}
			if err5 != nil {
				utils.SendResponse(w, 500, err5.Error(), "error", nil)
				return
			}

			context.Set(r, "currentUserAuthority", userStruct.Authority)
			context.Set(r, "currentUserEmail", userStruct.Email)
			context.Set(r, "currentUserKey", splitStrings[0])

			next.ServeHTTP(w, r)
		}
	})
}

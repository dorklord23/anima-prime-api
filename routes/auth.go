/* Copyright 2019 Tri Rumekso Anggie Wibowo (trirawibowo [at] gmail [dot] com)
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package routes

import (
	"encoding/json"
	"net/http"

	"github.com/dorklord23/anima-prime/models"
	"github.com/dorklord23/anima-prime/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// AuthenticateUser : endpoint to refresh access token and refresh token with login.
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	loginData := make(map[string]interface{})
	var user models.User
	var key *datastore.Key
	var err2 error
	requiredArgs := map[string]string{
		"Password": "required",
		"Email":    "required",
	}

	// Parse the request body and populate loginData
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(loginData, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Lookup the email
	q := datastore.NewQuery("users").Filter("Email =", loginData["Email"])
	t := q.Run(ctx)

	for {
		key, err2 = t.Next(&user)
		if err2 == datastore.Done {
			data := make(map[string]string)
			data["Message"] = "Wrong password or email"
			utils.SendResponse(w, 401, data, "fail", nil)
			return
		}
		if err2 != nil {
			utils.SendResponse(w, 500, err2.Error(), "error", nil)
			return
		}

		if checkPasswordHash(loginData["Password"].(string), user.Hash) {
			// Proceed to generate the access token and refresh token
			break
		} else {
			data := make(map[string]string)
			data["Message"] = "Wrong password or email"
			utils.SendResponse(w, 401, data, "fail", nil)
			return
		}
	}

	// Generate the access token and refresh token
	accessToken := utils.GenerateAccessToken(key)
	refreshToken := utils.RandSeq(20)

	// Update the user data because refresh token is stored in datastore
	user.RefreshToken = refreshToken

	// Commit to to server
	_, err3 := datastore.Put(ctx, key, &user)
	if err3 != nil {
		utils.SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	response := map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	}

	utils.SendResponse(w, 200, response, "success", nil)
}

// RefreshAccessToken : endpoint to refresh access token and refresh token without login.
func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var key *datastore.Key
	var err error
	ctx := appengine.NewContext(r)

	// Lookup the refresh token
	q := datastore.NewQuery("users").Filter("RefreshToken =", r.Header.Get("anima-prime-refresh-token"))
	t := q.Run(ctx)

	for {
		key, err = t.Next(&user)
		if err == datastore.Done {
			data := make(map[string]string)
			data["Message"] = "There is no such refresh token"
			utils.SendResponse(w, 404, data, "fail", nil)
			return
		}
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		// Proceed
		break
	}

	// Generate the access token and refresh token
	accessToken := utils.GenerateAccessToken(key)
	refreshToken := utils.RandSeq(20)

	// Update the user data because refresh token is stored in datastore
	user.RefreshToken = refreshToken

	// Commit to to server
	_, err3 := datastore.Put(ctx, key, &user)
	if err3 != nil {
		utils.SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	response := map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	}

	utils.SendResponse(w, 200, response, "success", nil)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

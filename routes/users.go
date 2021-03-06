/* Copyright 2019 Tri Rumekso Anggie Wibowo (trirawibowo [at] gmail [dot] com)
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dorklord23/anima-prime/models"
	"github.com/dorklord23/anima-prime/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// User : struct to hold user data to commit to Datastore
/* type User struct {
	FullName     string
	Email        string
	Hash         string
	Authority    string
	RefreshToken string
	CreatedAt    time.Time
	ModifiedAt   time.Time
} */

// CreateUsers : endpoint to create a new user and obtain access token
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	userMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"FullName":        "required",
		"Password":        "required",
		"PasswordConfirm": "required",
		"Email":           "required",
	}

	// Parse the request body and populate user
	err := json.NewDecoder(r.Body).Decode(&userMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(userMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Check if the password and the password confirm is exactly the same
	if userMap["Password"] != userMap["PasswordConfirm"] {
		data := make(map[string]string)
		data["PasswordConfirm"] = "Make sure this field is exactly the same with Password"
		utils.SendResponse(w, 400, data, "fail", nil)
		return
	}

	// Check if the email has been used already
	q := datastore.NewQuery("users").Filter("Email =", userMap["Email"])
	for t := q.Run(ctx); ; {
		var x models.User
		_, err := t.Next(&x)

		if err == datastore.Done {
			// No such email so it's safe to proceed
			break
		}
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		// The email has already been used
		data := make(map[string]string)
		data["Email"] = "The email has already been used"
		utils.SendResponse(w, 409, data, "fail", nil)
		return
	}

	// Hash the password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(userMap["Password"].(string)), 14)
	hash := utils.BytesToString(bytes)

	// Generate refresh token
	refreshToken := utils.RandSeq(20)

	// Preparing data to save
	userMap["Hash"] = hash
	userMap["RefreshToken"] = refreshToken
	userMap["CreatedAt"] = time.Now()
	userMap["ModifiedAt"] = time.Now()
	delete(userMap, "Password")
	delete(userMap, "PasswordConfirm")

	if userMap["Authority"] == nil {
		userMap["Authority"] = "regular"
	}

	var userStruct models.User
	err3 := mapstructure.Decode(userMap, &userStruct)
	if err3 != nil {
		utils.SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	// Save to Datastore
	userKey, err4 := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "users", nil), &userStruct)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	token := utils.GenerateAccessToken(userKey)

	data := make(map[string]string)
	options := make(map[string]string)
	location := fmt.Sprintf("%v://%v/api/users/%v", r.URL.Scheme, r.Host, userKey.Encode())
	options["Location"] = location
	data["Token"] = token
	data["RefreshToken"] = refreshToken

	utils.SendResponse(w, 201, data, "success", options)
}

// UpdateUsers : endpoint to update a user data
// A user could only change their own profile unless they're the admin
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"FullName": "optional",
		"Email":    "optional",
	}

	err := json.NewDecoder(r.Body).Decode(&userMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(userMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Check if this user is eligible to update the target profile by comparing the email in access token with the email of target profile IF AND ONLY IF the user is a non-admin

	key, err3 := datastore.DecodeKey(params["userKey"])
	if err3 != nil {
		data := make(map[string]string)
		data["Message"] = err3.Error()
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}

	// Because Datastore doesn't differentiate between creating and updating entity,
	// we need to retrieve the old data first and modify it before commiting it to Datastore
	var userStruct models.User

	// Retrieve the old data
	err5 := datastore.Get(ctx, key, &userStruct)
	if err5 == datastore.Done {
		// No such user
		data := make(map[string]string)
		data["Message"] = "There is no such user to update"
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}
	if err5 != nil {
		utils.SendResponse(w, 500, err5.Error(), "error", nil)
		return
	}

	// Check the requester's authority first
	currentUserAuthority := context.Get(r, "currentUserAuthority")
	if currentUserAuthority != utils.AdminAuthority {
		// Proceed to compare the emails
		currentUserEmail := context.Get(r, "currentUserEmail")
		if userStruct.Email != currentUserEmail {
			// Different email. Hence, the user is not eligible to update the target profile
			data := make(map[string]string)
			data["Message"] = "You are not eligible to update this user"
			utils.SendResponse(w, 403, data, "fail", nil)
			return
		}
	}

	// Overwrite it with the new one
	err2 := mapstructure.Decode(userMap, &userStruct)
	if err2 != nil {
		utils.SendResponse(w, 500, err2.Error(), "error", nil)
		return
	}

	// Commit it to Datastore
	_, err4 := datastore.Put(ctx, key, &userStruct)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	data["Message"] = "OK"

	utils.SendResponse(w, 204, data, "success", nil)
}

// GetUsers : endpoint to retrieve a user data
// A user could only retrieve their own profile unless they're the admin
func GetUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userMap := make(map[string]interface{})
	responseTemplate := make(map[string]string)
	ctx := appengine.NewContext(r)

	// Check if this user is eligible to retrieve the target profile by comparing the email in access token with the email of target profile IF AND ONLY IF the user is a non-admin
	key, err := datastore.DecodeKey(params["userKey"])
	if err != nil {
		data := make(map[string]string)
		data["Message"] = err.Error()
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}

	var userStruct models.User

	// Retrieve the data
	err2 := datastore.Get(ctx, key, &userStruct)
	if err2 == datastore.Done {
		// No such user
		data := make(map[string]string)
		data["Message"] = "There is no such user to retrieve"
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}
	if err2 != nil {
		utils.SendResponse(w, 500, err2.Error(), "error", nil)
		return
	}

	// Check the requester's authority first
	currentUserAuthority := context.Get(r, "currentUserAuthority")
	if currentUserAuthority != utils.AdminAuthority {
		// Proceed to compare the emails
		currentUserEmail := context.Get(r, "currentUserEmail")
		if userStruct.Email != currentUserEmail {
			// Different email. Hence, the user is not to retrieve the target profile
			data := make(map[string]string)
			data["Message"] = "You are not eligible to retrieve this user data"
			utils.SendResponse(w, 403, data, "fail", nil)
			return
		}
	}

	err3 := mapstructure.Decode(userStruct, &userMap)
	if err3 != nil {
		utils.SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	delete(userMap, "Hash")
	delete(userMap, "RefreshToken")
	delete(userMap, "CreatedAt")
	delete(userMap, "ModifiedAt")

	for key, value := range userMap {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		// var strValue string
		//
		// if strKey == "CreatedAt" || strKey == "ModifiedAt" {
		// 	strValue = fmt.Sprintf("%v", value.(time.Time).Format(time.RFC3339))
		// } else {
		// 	strValue = fmt.Sprintf("%v", value)
		// }

		responseTemplate[strKey] = strValue
	}

	utils.SendResponse(w, 200, responseTemplate, "success", nil)
}

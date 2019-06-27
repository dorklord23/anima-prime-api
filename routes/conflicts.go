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
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Conflict : data structure for conflicts
/* type Conflict struct {
	Name        string
	Description string
	Goal        string
	Difficulty  int
	Targets     []string
	IsResolved  bool
	ParentKey   string
} */

// CreateConflicts : endpoint to create a new conflict
func CreateConflicts(w http.ResponseWriter, r *http.Request) {
	requiredArgs := map[string]string{
		"Name":        "required",
		"Description": "required",
		"Goal":        "required",
		"Difficulty":  "required",
		"Targets":     "required",
	}

	conflictMap := make(map[string]interface{})
	conflictMap["IsResolved"] = false
	conflictMap["ParentKey"] = context.Get(r, "currentUserKey")

	var conflict models.Conflict
	models.CreateResource("conflicts", requiredArgs, conflictMap, conflict, w, r)
	/* conflictMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"Name":        "required",
		"Description": "required",
		"Goal":        "required",
		"Difficulty":  "required",
		"Targets":     "required",
	}

	// Parse the request body and populate user
	err := json.NewDecoder(r.Body).Decode(&conflictMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(conflictMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	var conflictStruct Conflict
	conflictMap["ParentKey"] = context.Get(r, "currentUserKey")
	conflictMap["IsResolved"] = false
	err3 := mapstructure.Decode(conflictMap, &conflictStruct)
	if err3 != nil {
		utils.SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	// Save to Datastore
	conflictKey, err4 := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "conflicts", nil), &conflictStruct)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	options := make(map[string]string)
	location := fmt.Sprintf("%v://%v/api/conflicts/%v", r.URL.Scheme, r.Host, conflictKey.Encode())
	data["ID"] = conflictKey.Encode()
	options["Location"] = location

	utils.SendResponse(w, 201, data, "success", options) */
}

// UpdateConflicts : endpoint to update a scene
func UpdateConflicts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	conflictMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"Name":        "optional",
		"Description": "optional",
		"Goal":        "optional",
		"Difficulty":  "optional",
		"Targets":     "optional",
		"IsResolved":  "optional",
	}

	err := json.NewDecoder(r.Body).Decode(&conflictMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(conflictMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Check if this user is authorized to update the target scene by comparing access token's user key with the parent key of target scene
	key, err3 := datastore.DecodeKey(params["conflictKey"])
	if err3 != nil {
		data := make(map[string]string)
		data["Message"] = err3.Error()
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}

	// Because Datastore doesn't differentiate between creating and updating entity,
	// we need to retrieve the old data first and modify it before commiting it to Datastore
	var conflict models.Conflict

	// Retrieve the old data
	err5 := datastore.Get(ctx, key, &conflict)
	if err5 == datastore.Done {
		// No such user
		data := make(map[string]string)
		data["Message"] = "There is no such conflict to update"
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
		// Proceed to compare the keys
		currentUserKey := context.Get(r, "currentUserKey")
		if conflict.ParentKey != currentUserKey {
			// Different key. Hence, the user is not authorized to update the target character
			data := make(map[string]string)
			data["Message"] = "You are not authorized to update this conflict"
			utils.SendResponse(w, 403, data, "fail", nil)
			return
		}
	}

	// Overwrite it with the new one
	err2 := mapstructure.Decode(conflictMap, &conflict)
	if err2 != nil {
		utils.SendResponse(w, 500, err2.Error(), "error", nil)
		return
	}

	// Commit it to Datastore
	_, err4 := datastore.Put(ctx, key, &conflict)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	data["Message"] = "OK"

	utils.SendResponse(w, 204, data, "success", nil)
}

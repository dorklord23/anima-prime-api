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

	"github.com/dorklord23/anima-prime/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Skill : data structure for character skills
type Skill struct {
	ID     string
	Rating int
}

// Trait : data structure for character traits
type Trait struct {
	Value    string
	IsTicked bool
}

// Character : data structure for characters
type Character struct {
	Name       string
	Concept    string
	Mark       string
	Passion    string
	Traits     []Trait
	Skills     []Skill
	Powers     []string
	Background string
	Links      []string
	ParentKey  string
}

// CreateCharacters : endpoint to create a new character (both PC and NPC)
func CreateCharacters(w http.ResponseWriter, r *http.Request) {
	characterMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"Name":       "required",
		"Concept":    "required",
		"Mark":       "required",
		"Passion":    "required",
		"Traits":     "required",
		"Skills":     "required",
		"Powers":     "required",
		"Background": "required",
		"Links":      "required",
	}

	// Parse the request body and populate user
	err := json.NewDecoder(r.Body).Decode(&characterMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(characterMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	var characterStruct Character
	characterMap["ParentKey"] = context.Get(r, "currentUserKey")
	err3 := mapstructure.Decode(characterMap, &characterStruct)
	if err3 != nil {
		utils.SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	// Save to Datastore
	characterKey, err4 := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "characters", nil), &characterStruct)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	options := make(map[string]string)
	location := fmt.Sprintf("%v://%v/api/characters/%v", r.URL.Scheme, r.Host, characterKey.Encode())
	data["ID"] = characterKey.Encode()
	options["Location"] = location

	utils.SendResponse(w, 201, data, "success", options)
}

// UpdateCharacters : endpoint to update a character (both PC and NPC)
func UpdateCharacters(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	characterMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"Name":       "optional",
		"Concept":    "optional",
		"Mark":       "optional",
		"Passion":    "optional",
		"Traits":     "optional",
		"Skills":     "optional",
		"Powers":     "optional",
		"Background": "optional",
		"Links":      "optional",
	}

	err := json.NewDecoder(r.Body).Decode(&characterMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(characterMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Check if this user is eligible to update the target character by comparing access token's user key with the parent key of target character

	key, err3 := datastore.DecodeKey(params["characterKey"])
	if err3 != nil {
		data := make(map[string]string)
		data["Message"] = err3.Error()
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}

	// Because Datastore doesn't differentiate between creating and updating entity,
	// we need to retrieve the old data first and modify it before commiting it to Datastore
	var character Character

	// Retrieve the old data
	err5 := datastore.Get(ctx, key, &character)
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
		currentUserKey := context.Get(r, "currentUserKey")
		if character.ParentKey != currentUserKey {
			// Different key. Hence, the user is not eligible to update the target character
			data := make(map[string]string)
			data["Message"] = "You are not eligible to update this character"
			utils.SendResponse(w, 403, data, "fail", nil)
			return
		}
	}

	// Overwrite it with the new one
	err2 := mapstructure.Decode(characterMap, &character)
	if err2 != nil {
		utils.SendResponse(w, 500, err2.Error(), "error", nil)
		return
	}

	// Commit it to Datastore
	_, err4 := datastore.Put(ctx, key, &character)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	data["Message"] = "OK"

	utils.SendResponse(w, 204, data, "success", nil)
}

// GetCharacters : endpoint to retrieve a character (both PC and NPC)
func GetCharacters(w http.ResponseWriter, r *http.Request) {}

// DeleteCharacters : endpoint to retrieve a character (both PC and NPC)
func DeleteCharacters(w http.ResponseWriter, r *http.Request) {}

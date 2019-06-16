/* Copyright 2019 Tri Rumekso Anggie Wibowo (trirawibowo [at] gmail [dot] com)
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package routes

import (
	"net/http"
	"strconv"

	"github.com/dorklord23/anima-prime/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// ChangeTraitTick : change a character's trait's tick status (tick/untick)
func ChangeTraitTick(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx := appengine.NewContext(r)
	// characterMap := make(map[string]interface{})

	// Check if this user is authorized to update the target character by comparing access token's user key with the parent key of target character

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
		data["Message"] = "There is no such character"
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
			// Different key. Hence, the user is not authorized to update the target character
			data := make(map[string]string)
			data["Message"] = "You are not authorized to update this character"
			utils.SendResponse(w, 403, data, "fail", nil)
			return
		}
	}

	// Change the tick status
	index, err := strconv.Atoi(params["traitIndex"])
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	characterTraits := make(map[int]Trait)

	for i := 0; i < len(character.Traits); i++ {
		characterTraits[i] = character.Traits[i]
	}

	if _, ok := characterTraits[index]; !ok {
		data := make(map[string]string)
		data["Message"] = "There is no trait with specified index"
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}

	character.Traits[index].IsTicked = !character.Traits[index].IsTicked
	/* characterMap["Traits"] = characterTraits

	// Overwrite it with the new one
	err2 := mapstructure.Decode(characterMap, &character)
	if err2 != nil {
		utils.SendResponse(w, 500, err2.Error(), "error", nil)
		return
	} */

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

// Reroll : reroll arbitrary number of dice
func Reroll(w http.ResponseWriter, r *http.Request) {
	dieQty, err := strconv.Atoi(r.URL.Query().Get("dieQty"))
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	data := make(map[string][]int)
	data["Dice"] = utils.DieRandomizer(dieQty)

	utils.SendResponse(w, 200, data, "success", nil)
}

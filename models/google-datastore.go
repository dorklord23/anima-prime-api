/* Copyright 2019 Tri Rumekso Anggie Wibowo (trirawibowo [at] gmail [dot] com)
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

// Not finished yet. Still figuring out how to implement separate function to handle saving data

package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/dorklord23/anima-prime/utils"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// CreateResource : function to handle creating new entity in Google Datastore
func CreateResource(resourceName string, requiredArgs map[string]string, resourceMap map[string]interface{}, resourceStruct interface{}, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var err error

	// Parse the request body and populate user
	err = json.NewDecoder(r.Body).Decode(&resourceMap)
	if err != nil {
		utils.SendResponse(w, 500, "1 "+err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(resourceMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Save to Datastore
	var resourceKey *datastore.Key
	// var err4 error

	switch resourceType := resourceStruct.(type) {
	// case Character, Conflict, Eidolon, Modifier, StatusChange, Power, SceneBonus, Scene, User:
	// Each case must uses one type only. Otherwise, it won't be saved properly (saved as new data but empty)
	case Character:
		err = mapstructure.Decode(resourceMap, &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		resourceKey, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, resourceName, nil), &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	case Conflict:
		err = mapstructure.Decode(resourceMap, &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		resourceKey, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, resourceName, nil), &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	case Eidolon:
		err = mapstructure.Decode(resourceMap, &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		resourceKey, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, resourceName, nil), &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	case Power:
		err = mapstructure.Decode(resourceMap, &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		resourceKey, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, resourceName, nil), &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	case Scene:
		err = mapstructure.Decode(resourceMap, &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		resourceKey, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, resourceName, nil), &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	case User:
		err = mapstructure.Decode(resourceMap, &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}

		resourceKey, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, resourceName, nil), &resourceType)
		if err != nil {
			utils.SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	default:
		text := reflect.TypeOf(resourceStruct).Name() // Kind().String()
		utils.SendResponse(w, 500, "Unknown resource type: "+text, "error", nil)
		return
	}

	data := make(map[string]string)
	options := make(map[string]string)
	location := fmt.Sprintf("%v://%v/api/%v/%v", r.URL.Scheme, r.Host, resourceName, resourceKey.Encode())
	data["ID"] = resourceKey.Encode()
	options["Location"] = location

	utils.SendResponse(w, 201, data, "success", options)
}

// UpdateResource : function to handle updating entity in Google Datastore
func UpdateResource(
	requiredArgs map[string]string,
	resourceStruct interface{},
	params map[string]string,
	w http.ResponseWriter,
	r *http.Request) {
	resourceMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)

	err := json.NewDecoder(r.Body).Decode(&resourceMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(resourceMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Check if this user is authorized to update the target scene by comparing access token's user key with the parent key of target scene
	key, err3 := datastore.DecodeKey(params["resourceKey"])
	if err3 != nil {
		data := make(map[string]string)
		data["Message"] = err3.Error()
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}

	// Because Datastore doesn't differentiate between creating and updating entity,
	// we need to retrieve the old data first and modify it before commiting it to Datastore

	// Retrieve the old data
	err5 := datastore.Get(ctx, key, &resourceStruct)
	if err5 == datastore.Done {
		// No such user
		data := make(map[string]string)
		data["Message"] = "There is no such resource to update"
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
		// TODO: change resourceStruct to map
		if resourceStruct.(map[string]interface{})["ParentKey"] != currentUserKey {
			// Different key. Hence, the user is not authorized to update the target character
			data := make(map[string]string)
			data["Message"] = "You are not authorized to update this resource"
			utils.SendResponse(w, 403, data, "fail", nil)
			return
		}
	}

	// Overwrite it with the new one
	err2 := mapstructure.Decode(resourceMap, &resourceStruct)
	if err2 != nil {
		utils.SendResponse(w, 500, err2.Error(), "error", nil)
		return
	}

	// Commit it to Datastore
	_, err4 := datastore.Put(ctx, key, &resourceStruct)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	data["Message"] = "OK"

	utils.SendResponse(w, 204, data, "success", nil)
}

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

// SceneBonus : data structure for scene bonus
type SceneBonus struct {
	BonusID string
	UserID  string
}

// Scene : data structure for scenes
type Scene struct {
	Name        string
	Description string
	IsResolved  bool
	Bonus       []SceneBonus
	ParentKey   string
}

// CreateScenes : endpoint to create a new scene
func CreateScenes(w http.ResponseWriter, r *http.Request) {
	sceneMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"Name":        "required",
		"Description": "required",
	}

	// Parse the request body and populate user
	err := json.NewDecoder(r.Body).Decode(&sceneMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(sceneMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	var sceneStruct Scene
	sceneMap["ParentKey"] = context.Get(r, "currentUserKey")
	sceneMap["IsResolved"] = false
	err3 := mapstructure.Decode(sceneMap, &sceneStruct)
	if err3 != nil {
		utils.SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	// Save to Datastore
	sceneKey, err4 := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "scenes", nil), &sceneStruct)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	options := make(map[string]string)
	location := fmt.Sprintf("%v://%v/api/scenes/%v", r.URL.Scheme, r.Host, sceneKey.Encode())
	data["ID"] = sceneKey.Encode()
	options["Location"] = location

	utils.SendResponse(w, 201, data, "success", options)
}

// UpdateScenes : endpoint to update a scene
func UpdateScenes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sceneMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	requiredArgs := map[string]string{
		"Name":        "optional",
		"Description": "optional",
		"IsResolved":  "optional",
		"Bonus":       "optional",
	}

	err := json.NewDecoder(r.Body).Decode(&sceneMap)
	if err != nil {
		utils.SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := utils.CheckArgs(sceneMap, requiredArgs)
	if missingArgs != nil {
		utils.SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	// Check if this user is authorized to update the target scene by comparing access token's user key with the parent key of target scene
	key, err3 := datastore.DecodeKey(params["sceneKey"])
	if err3 != nil {
		data := make(map[string]string)
		data["Message"] = err3.Error()
		utils.SendResponse(w, 404, data, "fail", nil)
		return
	}

	// Because Datastore doesn't differentiate between creating and updating entity,
	// we need to retrieve the old data first and modify it before commiting it to Datastore
	var scene Scene

	// Retrieve the old data
	err5 := datastore.Get(ctx, key, &scene)
	if err5 == datastore.Done {
		// No such user
		data := make(map[string]string)
		data["Message"] = "There is no such scene to update"
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
		if scene.ParentKey != currentUserKey {
			// Different key. Hence, the user is not authorized to update the target character
			data := make(map[string]string)
			data["Message"] = "You are not authorized to update this character"
			utils.SendResponse(w, 403, data, "fail", nil)
			return
		}
	}

	// Overwrite it with the new one
	err2 := mapstructure.Decode(sceneMap, &scene)
	if err2 != nil {
		utils.SendResponse(w, 500, err2.Error(), "error", nil)
		return
	}

	// Commit it to Datastore
	_, err4 := datastore.Put(ctx, key, &scene)
	if err4 != nil {
		utils.SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	data["Message"] = "OK"

	utils.SendResponse(w, 204, data, "success", nil)
}

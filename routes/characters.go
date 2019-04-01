package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Skill : data structure for character skills
type Skill struct {
	ID     string
	Rating int
}

// Character : data structure for characters
type Character struct {
	Name       string
	Concept    string
	Mark       string
	Passion    string
	Traits     []string
	Skills     []Skill
	Powers     []string
	Background string
	Links      []string
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
		SendResponse(w, 500, err.Error(), "error", nil)
		return
	}

	// Check if there are any missing arguments
	missingArgs := CheckArgs(characterMap, requiredArgs)
	if missingArgs != nil {
		SendResponse(w, 400, missingArgs, "fail", nil)
		return
	}

	var characterStruct Character
	err3 := mapstructure.Decode(characterMap, &characterStruct)
	if err3 != nil {
		SendResponse(w, 500, err3.Error(), "error", nil)
		return
	}

	// Save to Datastore
	characterKey, err4 := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "characters", nil), &characterStruct)
	if err4 != nil {
		SendResponse(w, 500, err4.Error(), "error", nil)
		return
	}

	data := make(map[string]string)
	options := make(map[string]string)
	location := fmt.Sprintf("%v://%v/api/characters/%v", r.URL.Scheme, r.Host, characterKey.Encode())
	data["ID"] = characterKey.Encode()
	options["Location"] = location

	SendResponse(w, 201, data, "success", options)
}

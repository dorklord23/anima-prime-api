package routes

import (
	"github.com/dorklord23/anima-prime/utils"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type user struct {
	FullName   string
	Email      string
	Hash       string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// CreateUsers : endpoint to create a new user and obtain access token
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	userMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)
	var template strings.Builder

	// Parse the request body and populate user
	err := json.NewDecoder(r.Body).Decode(&userMap)
	if err != nil {
		SendResponse(w, 500, err.Error(), "error")
		return
	}

	// Check if the password and the password confirm is exactly the same
	if userMap["Password"] != userMap["PasswordConfirm"] {
		data := make(map[string]string)
		data["PasswordConfirm"] = "Make sure this field is exactly the same with Password"
		SendResponse(w, 400, data, "fail")
		return
	}

	// Hash the password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(userMap["Password"].(string)), 14)
	hash := utils.BytesToString(bytes)

	// Generate the access token
	template.WriteString(userMap["Email"].(string))
	template.WriteString("|")
	template.WriteString(time.Now().Format(time.RFC3339))
	template.WriteString("|")
	// Expiry time for the token in days
	template.WriteString("90")
	token := base64.StdEncoding.EncodeToString([]byte(template.String()))

	// Preparing data to save
	userMap["Hash"] = hash
	userMap["CreatedAt"] = time.Now()
	userMap["ModifiedAt"] = time.Now()
	delete(userMap, "Password")
	delete(userMap, "PasswordConfirm")

	var userStruct user
	err3 := mapstructure.Decode(userMap, &userStruct)
	if err3 != nil {
		SendResponse(w, 500, err3.Error(), "error")
		return
	}

	// Save to Datastore
	_, err4 := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "users", nil), &userStruct)
	if err4 != nil {
		SendResponse(w, 500, err4.Error(), "error")
		return
	}

	data := make(map[string]string)
	data["Token"] = token

	SendResponse(w, 200, data, "success")
}

// UpdateUsers : endpoint to update a user.
// A user could only change their own profile unless they're the admin
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userMap := make(map[string]interface{})
	ctx := appengine.NewContext(r)

	err := json.NewDecoder(r.Body).Decode(&userMap)
	if err != nil {
		SendResponse(w, 500, err.Error(), "error")
		return
	}

	// TODO: A user could only change their own profile unless they're the admin

	var userStruct user
	err2 := mapstructure.Decode(userMap, &userStruct)
	if err2 != nil {
		SendResponse(w, 500, err2.Error(), "error")
		return
	}

	// params["userId"]
}

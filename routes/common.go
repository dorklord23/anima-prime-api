package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dorklord23/anima-prime/utils"
	"google.golang.org/appengine/datastore"
)

type errorResponse struct {
	Status  string
	Message string
}

type successResponse struct {
	Status string
	Data   map[string]string
}

type failResponse struct {
	Status string
	Data   map[string]string
}

// SendResponse : function to send JSON-formatted HTTP response
func SendResponse(w http.ResponseWriter, statusCode int, payload interface{}, responseType string, options map[string]string) {
	var response []byte
	var err error

	switch responseType {
	case "success":
		template := successResponse{
			Status: "success",
			Data:   payload.(map[string]string),
		}

		response, err = json.Marshal(template)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	case "fail":
		template := failResponse{
			Status: "fail",
			Data:   payload.(map[string]string),
		}

		response, err = json.Marshal(template)
		if err != nil {
			SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	case "error":
		template := errorResponse{
			Status:  "error",
			Message: payload.(string),
		}

		response, err = json.Marshal(template)
		if err != nil {
			SendResponse(w, 500, err.Error(), "error", nil)
			return
		}
	default:
	}

	if options != nil {
		for key, value := range options {
			w.Header().Set(key, value)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// GenerateAccessToken : function to generate access token
func GenerateAccessToken(userKey *datastore.Key) string {
	var template strings.Builder

	// Generate the access token
	template.WriteString(userKey.Encode())
	template.WriteString("|")
	template.WriteString(time.Now().Format(time.RFC3339))
	template.WriteString("|")
	// Expiry time for the token in 1 days
	template.WriteString("1")
	template.WriteString("|")
	template.WriteString(utils.RandSeq(5))

	return base64.StdEncoding.EncodeToString([]byte(template.String()))
}

// CheckArgs : function to check if an endpoint's arguments are already supplied completely
func CheckArgs(suppliedArgs map[string]interface{}, requiredArgs map[string]string) map[string]string {
	var errorList []string

	for key, value := range requiredArgs {
		if suppliedArgs[key] == nil && value == "required" {
			errorList = append(errorList, key)
		}
	}

	if len(errorList) > 0 {
		// Missing arguments
		result := make(map[string]string)

		for _, value := range errorList {
			result[value] = "This argument is missing from the request"
		}

		return result
	}

	// Everything is fine
	return nil
}

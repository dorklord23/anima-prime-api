package routes

import (
	"encoding/json"
	"net/http"
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

// SendResponse is a function to send JSON-formatted HTTP response
func SendResponse(w http.ResponseWriter, statusCode int, payload interface{}, responseType string) {
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
			SendResponse(w, 500, err.Error(), "error")
			return
		}
	case "fail":
		template := failResponse{
			Status: "fail",
			Data:   payload.(map[string]string),
		}

		response, err = json.Marshal(template)
		if err != nil {
			SendResponse(w, 500, err.Error(), "error")
			return
		}
	case "error":
		template := errorResponse{
			Status:  "error",
			Message: payload.(string),
		}

		response, err = json.Marshal(template)
		if err != nil {
			SendResponse(w, 500, err.Error(), "error")
			return
		}
	default:
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

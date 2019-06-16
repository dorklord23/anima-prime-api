/* Copyright 2019 Tri Rumekso Anggie Wibowo (trirawibowo [at] gmail [dot] com)
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package utils

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
)

// SendResponse : function to send JSON-formatted HTTP response
func SendResponse(w http.ResponseWriter, statusCode int, payload interface{}, responseType string, options map[string]string) {
	template := make(map[string]interface{})
	template["Status"] = responseType

	if responseType == "error" {
		template["Message"] = payload
	} else {
		template["Data"] = payload
	}

	response, err := json.Marshal(template)
	if err != nil {
		SendResponse(w, 500, err.Error(), "error", nil)
		return
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
	template.WriteString(RandSeq(5))

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

// DieRandomizer : randomizes arbitrary number of dice roll
func DieRandomizer(dieQty int) []int {
	result := make([]int, dieQty)
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 6

	for i := 0; i < dieQty; i++ {
		result[i] = rand.Intn(max-min) + min
	}

	return result
}

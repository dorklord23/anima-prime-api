// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"

	"github.com/dorklord23/anima-prime/middlewares"
	"github.com/dorklord23/anima-prime/routes"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func init() {
	// A little hack to use mux in App Engine
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	// The only non-RESTful endpoint in this API to accomodate login
	s.HandleFunc("/login", routes.AuthenticateUser).Methods("POST")

	s.HandleFunc("/users", routes.CreateUsers).Methods("POST")
	s.HandleFunc("/users/{userKey}", routes.UpdateUsers).Methods("PUT")
	s.HandleFunc("/users/{userKey}", routes.GetUsers).Methods("GET")

	s.HandleFunc("/characters", routes.CreateCharacters).Methods("POST")
	s.HandleFunc("/characters/{characterKey}", routes.UpdateCharacters).Methods("PUT")
	s.HandleFunc("/characters/{characterKey}", routes.GetCharacters).Methods("GET")
	s.HandleFunc("/characters/{characterKey}", routes.DeleteCharacters).Methods("DELETE")

	s.HandleFunc("/tokens", routes.RefreshAccessToken).Methods("GET")

	s.Use(middlewares.Authenticate)
	// The path "/" matches everything not matched by some other path.
	http.Handle("/", r)
}

func main() {
	appengine.Main()
}

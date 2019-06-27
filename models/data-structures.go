/* Copyright 2019 Tri Rumekso Anggie Wibowo (trirawibowo [at] gmail [dot] com)
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package models

import "time"

// characters.go

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

// conflicts.go

// Conflict : data structure for conflicts
type Conflict struct {
	Name        string
	Description string
	Goal        string
	Difficulty  int
	Targets     []string
	IsResolved  bool
	ParentKey   string
}

// eidolons.go

// Eidolon : data structure for eidolons
type Eidolon struct {
	Name        string
	Description string
	Level       int
	Type        int
	Skill       int
	Powers      []int
	Weakness    int
	ParentKey   string
}

// powers.go

// Modifier : data structure for modifiers
type Modifier struct {
	TargetProp  string
	Value       int
	IsPermanent bool
}

// StatusChange : a generic term for buffs and debuffs
type StatusChange struct {
	Name        string
	Description string
	Changes     []Modifier
}

// Power : data structure for powers
type Power struct {
	Name        string
	Description string
	Type        int
	// Effect : array of StatusChange keys
	Effect    []string
	ParentKey string
}

// ApplyEffect : apply this particular power's effect
func (p Power) ApplyEffect() {
	// Do nothing for now
}

// scenes.go

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

// users.go

// User : struct to hold user data to commit to Datastore
type User struct {
	FullName     string
	Email        string
	Hash         string
	Authority    string
	RefreshToken string
	CreatedAt    time.Time
	ModifiedAt   time.Time
}

package routes

import "net/http"

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
func CreateCharacters(w http.ResponseWriter, r *http.Request) {}

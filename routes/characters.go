package routes

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
	PassionID  int
	Traits     []string
	Skills     []Skill
	Powers     []string
	Background string
	Links      []string
}

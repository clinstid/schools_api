package db

import (
	"fmt"
)

// SchoolsResult is a struct used to represent a list of schools from the
// database that includes a slice of school names and the total count of
// schools in the database.
type SchoolsResult struct {
	Schools []string
	Total   int
}

// GetSchools returns a SchoolsResult struct based on the specified offset and
// limit into the slice of schools in the database.
func GetSchools(limit, offset int) SchoolsResult {
	total := len(schoolDB)
	result := SchoolsResult{Total: total}
	if offset < len(schoolDB)-1 {
		if offset+limit > total {
			limit = total - offset
		}
		result.Schools = schoolDB[offset : offset+limit]
	}
	return result
}

// GetSchool returns the name of a school specified by an id. The id is an
// index in to the slice of schools. If the school is not found nil and an
// error will be returned.
func GetSchool(id int) (*string, error) {
	if id < len(schoolDB) {
		name := schoolDB[id]
		return &name, nil
	} else {
		return nil, fmt.Errorf("School with id %d not found", id)
	}
}

// AddSchool adds a new school with the specified name to the end of the slice
// of schools.
func AddSchool(name string) int {
	schoolDB = append(schoolDB, name)
	return len(schoolDB) - 1
}

// UpdateSchool updates the name of a school at the specified id and returns
// the id of the school. If there is no school at the specified id then nil and
// error will be returned.
func UpdateSchool(id int, name string) (*int, error) {
	if id < len(schoolDB) {
		schoolDB[id] = name
		return &id, nil
	} else {
		return nil, fmt.Errorf("School with id %d not found", id)
	}
}

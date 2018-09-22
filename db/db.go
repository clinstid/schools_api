package db

import (
	"fmt"
)

type SchoolsResult struct {
	Schools []string
	Total   int
}

func GetSchools(limit, offset int) SchoolsResult {
	result := SchoolsResult{Total: len(schoolDB)}
	if offset < len(schoolDB)-1 {
		result.Schools = schoolDB[offset : offset+limit]
	}
	return result
}

func GetSchool(id int) (*string, error) {
	if id < len(schoolDB) {
		name := schoolDB[id]
		return &name, nil
	} else {
		return nil, fmt.Errorf("School with id %d not found", id)
	}
}

func AddSchool(name string) int {
	schoolDB = append(schoolDB, name)
	return len(schoolDB) - 1
}

func UpdateSchool(id int, name string) (*int, error) {
	if id < len(schoolDB) {
		schoolDB[id] = name
		return &id, nil
	} else {
		return nil, fmt.Errorf("School with id %d not found", id)
	}
}

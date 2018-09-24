package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clinstid/schools_api/db"
	"github.com/clinstid/schools_api/resources"
	"github.com/gin-gonic/gin"
)

const (
	// Query parameter constants
	limitField   = "limit"
	limitDefault = 100
	maxLimit     = 100
	minLimit     = 1

	offsetField   = "offset"
	offsetDefault = 0
	minOffSet     = 0
)

var (
	// Error messages
	limitNotNumberErrMsg    = "limit query parameter must be a number"
	limitOutOfBoundsErrMsg  = fmt.Sprintf("limit query parameter must be at least %d and no greater than %d", minLimit, maxLimit)
	offsetNotNumberErrMsg   = "offset query parameter must be a number"
	offsetOutOfBoundsErrMsg = fmt.Sprintf("offset query parameter must be at least %d", minOffSet)
	schoolIdNotNumberErrMsg = "school id must be a number"
)

// buildErrorResponse returns a gin.H struct with a message property that will
// look like the following when renderd as JSON:
//
// {
//   "message": "Error message"
// }
func buildErrorResponse(message string) gin.H {
	return gin.H{"message": message}
}

// buildBindErrorResponse builds a custom error response with the error type
// returned from a Bind call used to deserialize a request body into an
// internal structure.
func buildBindErrorResponse(err error) gin.H {
	switch actualErr := err.(type) {
	case *json.UnmarshalTypeError:
		msg := fmt.Sprintf("Field %q must be a %s", actualErr.Field, actualErr.Type)
		return buildErrorResponse(msg)
	default:
		return buildErrorResponse(err.Error())
	}

}

// buildSchoolLink returns a URL for the current school
func buildSchoolLink(r *http.Request, schoolID int) string {
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	host := r.Host
	path := fmt.Sprintf("%s/%d", r.URL.Path, schoolID)

	link := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	return link.String()

}

// buildListSchoolsLink returns a URL for making a ListSchools request with the
// offset and limit encoded as query parameters. This function is used for
// building the next, prev, first, and last links that are returned in a
// ListSchools response.
func buildListSchoolsLink(r *http.Request, offset int, limit int) string {
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	host := r.Host
	path := r.URL.Path

	link := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	q := link.Query()
	q.Set(offsetField, strconv.Itoa(offset))
	q.Set(limitField, strconv.Itoa(limit))
	link.RawQuery = q.Encode()
	return link.String()
}

// ListSchools is a handler function for for the ListSchools operation. It
// takes query parameters `limit` and `offset` to determine what schools in the
// list to return. The format of the response looks like:
//
// ```json
// {
//   "schools": [
//     {
//       "name": "School name",
//       "id": 1234
//     },
//     ...
//   ],
//   "meta": {
//		"total": 4567
//   },
//   "links": {
//		"first": "http://...",
//		"last": "http://...",
//		"next": "http://...",
//		"prev": "http://...",
//   }
// }
// ```
//
// `schools` is an array of school objects with a name and an id
// `meta` contains meta data about the collection including the total number of schools
// `links` has URLs for first, last, next, and previous pages of schools
func ListSchools(c *gin.Context) {
	// Parse query parameters
	limit, err := strconv.Atoi(c.DefaultQuery(limitField, strconv.Itoa(limitDefault)))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse(limitNotNumberErrMsg))
		return
	}

	if limit < minLimit || limit > maxLimit {
		c.JSON(http.StatusBadRequest, buildErrorResponse(limitOutOfBoundsErrMsg))
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery(offsetField, strconv.Itoa(offsetDefault)))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse(offsetNotNumberErrMsg))
		return
	}

	if offset < 0 {
		c.JSON(http.StatusBadRequest, buildErrorResponse(offsetOutOfBoundsErrMsg))
		return
	}

	// Retrieve the slice of schools for the given limit and offset.
	sResult := db.GetSchools(limit, offset)

	// Build links for pagination
	firstLink := buildListSchoolsLink(c.Request, 0, limit)

	total := sResult.Total
	pageCount := total / limit
	lastLink := buildListSchoolsLink(c.Request, limit*pageCount, limit)

	var nextLink string
	if offset+limit < total {
		nextLink = buildListSchoolsLink(c.Request, offset+limit, limit)
	}

	var prevLink string
	if offset-limit >= 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		} else if prevOffset > total {
			prevOffset = limit * pageCount
		}
		prevLink = buildListSchoolsLink(c.Request, prevOffset, limit)
	}

	// Build the response object
	schools := resources.Schools{
		Schools: make([]resources.School, 0, len(sResult.Schools)),
		Meta:    resources.Meta{Total: sResult.Total},
		Links: resources.Links{
			First: firstLink,
			Last:  lastLink,
			Next:  nextLink,
			Prev:  prevLink,
		},
	}
	// Add the schools to the response object
	for idx, s := range sResult.Schools {
		schools.Schools = append(schools.Schools, resources.School{ID: offset + idx, Name: s})
	}

	// Render and return the response
	c.JSON(http.StatusOK, schools)
}

// AddSchool adds a new school to the list. It takes a new school object in the
// request body:
//
// ```json
// {
// 	 "name": "New School Name"
// }```
func AddSchool(c *gin.Context) {
	var school resources.School
	err := c.ShouldBind(&school)
	if err != nil {
		c.JSON(http.StatusBadRequest, buildBindErrorResponse(err))
		return
	}

	schoolID := db.AddSchool(school.Name)
	c.Header("Location", buildSchoolLink(c.Request, schoolID))
	c.JSON(http.StatusCreated, resources.School{ID: schoolID, Name: school.Name})
}

// GetSchool retrieves a single school with the specified id.
func GetSchool(c *gin.Context) {
	// Get the id from the path
	schoolID, err := strconv.Atoi(c.Param("schoolID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse(schoolIdNotNumberErrMsg))
		return
	}

	// Look up the school in the database
	schoolName, err := db.GetSchool(schoolID)
	if err != nil {
		c.JSON(http.StatusNotFound, buildErrorResponse(err.Error()))
		return
	}

	// Build the response object
	school := resources.School{ID: schoolID, Name: *schoolName}

	// Render the response object
	c.JSON(http.StatusOK, school)
}

// UpdateSchool updates a single school with the specified id. The body
// contains the new name for the school.
func UpdateSchool(c *gin.Context) {
	// Get the id from the path
	schoolID, err := strconv.Atoi(c.Param("schoolID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse(schoolIdNotNumberErrMsg))
		return
	}

	// Deserialize the JSON body and bind it to the resources.School struct
	var school resources.School
	err = c.ShouldBind(&school)
	if err != nil {
		c.JSON(http.StatusBadRequest, buildBindErrorResponse(err))
		return
	}

	// Update the school in the database
	_, err = db.UpdateSchool(schoolID, school.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, buildErrorResponse(err.Error()))
		return
	}

	// Render the response
	c.JSON(http.StatusOK, resources.School{ID: schoolID, Name: school.Name})
}

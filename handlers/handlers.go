package handlers

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/clinstid/schools_api/db"
	"gitlab.com/clinstid/schools_api/resources"
)

const (
	limitField    = "limit"
	limitDefault  = "100"
	offsetField   = "offset"
	offsetDefault = "0"
)

func buildErrorResponse(message string) gin.H {
	return gin.H{"message": message}
}

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

func ListSchools(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery(limitField, limitDefault))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse("limit query parameter must be a number"))
		return
	}

	if limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, buildErrorResponse("limit query parameter must be at least one 1 and no greater than 100"))
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery(offsetField, offsetDefault))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse("offset query parameter must be a number"))
		return
	}

	if offset < 0 {
		c.JSON(http.StatusBadRequest, buildErrorResponse("limit query parameter must be at least 0"))
		return
	}

	sResult := db.GetSchools(limit, offset)

	firstLink := buildListSchoolsLink(c.Request, 0, limit)

	total := sResult.Total
	pageCount := total / limit
	lastLink := buildListSchoolsLink(c.Request, limit*(pageCount-1), limit)

	var nextLink string
	if offset+limit < total {
		nextLink = buildListSchoolsLink(c.Request, offset+limit, limit)
	}

	var prevLink string
	if offset > 0 && (offset < total-1) {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		prevLink = buildListSchoolsLink(c.Request, prevOffset, limit)
	}

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
	for idx, s := range sResult.Schools {
		schools.Schools = append(schools.Schools, resources.School{ID: offset + idx, Name: s})
	}

	c.JSON(http.StatusOK, schools)
}

func AddSchool(c *gin.Context) {
	var school resources.School
	err := c.ShouldBindJSON(&school)
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse(err.Error()))
	}

	schoolID := db.AddSchool(school.Name)
	c.JSON(http.StatusOK, resources.School{ID: schoolID, Name: school.Name})
}

func GetSchool(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse("school id must be a number"))
		return
	}

	schoolName, err := db.GetSchool(schoolID)
	if err != nil {
		c.JSON(http.StatusNotFound, buildErrorResponse(err.Error()))
		return
	}

	school := resources.School{ID: schoolID, Name: *schoolName}
	c.JSON(http.StatusOK, school)
}

func UpdateSchool(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse(err.Error()))
	}

	var school resources.School
	err = c.ShouldBindJSON(&school)
	if err != nil {
		c.JSON(http.StatusBadRequest, buildErrorResponse(err.Error()))
	}

	_, err = db.UpdateSchool(schoolID, school.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, buildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resources.School{ID: schoolID, Name: school.Name})
}

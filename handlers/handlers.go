package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/clinstid/schools_api/db"
	"gitlab.com/clinstid/schools_api/resources"
)

func ListSchools(c *gin.Context) {
	schoolList := db.GetSchools()

	schools := resources.Schools{
		Schools: make([]resources.School, 0, len(schoolList)),
	}
	for idx, s := range schoolList {
		schools.Schools = append(schools.Schools, resources.School{ID: idx, Name: s})
	}

	c.JSON(http.StatusOK, schools)
}

func AddSchool(c *gin.Context) {
	var school resources.School
	err := c.ShouldBindJSON(&school)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	schoolID := db.AddSchool(school.Name)
	c.JSON(http.StatusOK, resources.School{ID: schoolID, Name: school.Name})
}

func GetSchool(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	schoolName, err := db.GetSchool(schoolID)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	school := resources.School{ID: schoolID, Name: *schoolName}
	c.JSON(http.StatusOK, school)
}

func UpdateSchool(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	var school resources.School
	err = c.ShouldBindJSON(&school)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	_, err = db.UpdateSchool(schoolID, school.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resources.School{ID: schoolID, Name: school.Name})
}

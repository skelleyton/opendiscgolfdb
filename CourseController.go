package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Database *DB
}

func (con *Controller) GetCourse(c *gin.Context) {
	courseId := c.Param("courseId")

	course, err := con.Database.GetCourseById(&courseId)

	if err != nil {
		c.String(http.StatusNotFound, "Not Found")
	}

	c.JSON(200, course)
}

func (con *Controller) ListCourses(c *gin.Context) {
	type NextKey struct {
		ViewedRecords int `json:"viewedRecords"`
	}

	type Response struct {
		Courses []Course `json:"courses"`
		Count   int      `json:"count"`
		NextKey string   `json:"nextKey"`
	}

	base64Encoder := base64.StdEncoding.WithPadding(-1)

	perPage := c.Query("perpage")
	nextKeyStr := c.Query("nextkey")

	perPageInt, err := strconv.ParseInt(perPage, 10, 0)

	if err != nil {
		perPageInt = 10
	}

	var nextKey *NextKey
	viewedRecords := 0

	if nextKeyStr != "" {
		nextKeyByte := make([]byte, base64.StdEncoding.WithPadding(-1).DecodedLen(len(nextKeyStr)))

		readLen, err := base64Encoder.Decode(nextKeyByte, []byte(nextKeyStr))

		if err != nil {
			log.Fatal(err)
		}

		nextKeySlice := nextKeyByte[:readLen]

		if err := json.Unmarshal(nextKeySlice, &nextKey); err != nil {
			log.Print("Failed to Unmarshal nextkey")
			log.Fatal(err)
		}

		viewedRecords = nextKey.ViewedRecords
	}

	courses := *con.Database.ListCourses()

	coursesLen := len(courses)

	coursesSlice := courses[viewedRecords:]

	var responseNextKey NextKey

	if len(coursesSlice) > int(perPageInt) {
		coursesSlice = coursesSlice[:perPageInt]
		responseNextKey = NextKey{
			viewedRecords + int(perPageInt),
		}
	}

	var encodedResponseNextKey string

	if (responseNextKey != NextKey{}) {
		data, err := json.Marshal(responseNextKey)

		if err != nil {
			log.Fatal(err)
		}
		encodedResponseNextKey = base64Encoder.EncodeToString(data)
	}

	c.JSON(200, Response{
		coursesSlice,
		coursesLen,
		encodedResponseNextKey,
	})
}

func (con *Controller) Search(c *gin.Context) {
	boundingBoxString := c.Query("bounding_box")

	if boundingBoxString == "" {
		c.JSON(400, "Invalid Input")
		return
	}

	var boundingBox BoundingBox

	if err := json.Unmarshal([]byte(boundingBoxString), &boundingBox); err != nil {
		c.JSON(400, "invalid_bounding_box")
		return
	}

	if courses, err := con.Database.SearchCourses(boundingBox); err != nil {
		c.JSON(400, err.(error).Error())
	} else {
		c.JSON(200, courses)
	}
}

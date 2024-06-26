package databases

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type CourseDB struct {
	db *[]Course
}

func NewCourseDB(path string) *CourseDB {
	dbPath := "./db.json"

	if path != "" {
		dbPath = path
	}
	dbByte, err := os.ReadFile(dbPath)

	if err != nil {
		log.Fatal(err)
	}

	var db *[]Course

	if err := json.Unmarshal(dbByte, &db); err != nil {
		log.Fatal(err)
	}

	return &CourseDB{db}
}

func (d *CourseDB) GetCourseById(id *string) (Course, error) {
	for _, value := range *d.db {
		if value.Properties.ID == *id {
			return value, nil
		}
	}

	return Course{}, errors.New("CourseNotExist")
}

func (d *CourseDB) ListCourses() *[]Course {
	return d.db
}

func (d *CourseDB) SearchCoursesByBoundingBox(boundingBox BoundingBox) (*[]Course, error) {
	if (boundingBox != BoundingBox{}) {
		boundingBox, err := mapBoundingBox(boundingBox)

		if err != nil {
			return nil, err
		}

		var courses []Course

		for _, val := range *d.db {
			courseCoords := val.Geometry.Coordinates

			if boundingBox[0][0] < courseCoords[0] &&
				boundingBox[1][0] > courseCoords[0] &&
				boundingBox[0][1] < courseCoords[1] &&
				boundingBox[1][1] > courseCoords[1] {
				courses = append(courses, val)
			}

		}
		return &courses, nil
	}

	return nil, errors.New("invalid_search_param")
}

func (d *CourseDB) SearchCoursesByPostCode(postCode string) *[]Course {
	var courses []Course

	for _, val := range *d.db {
		if val.Properties.ZipCode == postCode {
			courses = append(courses, val)
		}
	}

	return &courses
}

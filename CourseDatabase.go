package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

type CourseProperties struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Course struct {
	Type       string           `json:"type"`
	Geometry   Geometry         `json:"geometry"`
	Properties CourseProperties `json:"properties"`
}

type DB struct {
	DB *[]Course
}

func (d *DB) GetCourseById(id *string) (Course, error) {
	for _, value := range *d.DB {
		if value.Properties.ID == *id {
			return value, nil
		}
	}

	return Course{}, errors.New("CourseNotExist")
}

func (d *DB) ListCourses() *[]Course {
	return d.DB
}

func NewDB(path string) *DB {
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

	return &DB{db}
}

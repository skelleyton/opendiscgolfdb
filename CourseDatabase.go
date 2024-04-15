package main

import (
	"errors"
)

type GeometryType string

const (
	Point GeometryType = "Point"
	Line  GeometryType = "Line"
)

type Geometry struct {
	Type        GeometryType `json:"type"`
	Coordinates []float32    `json:"coordinates"`
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

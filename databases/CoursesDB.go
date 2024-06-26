package databases

import (
	"errors"
	"log"
	"time"

	"opendgdb/types"

	"github.com/couchbase/gocb/v2"
)

// Coordinates consists of 2 floats in an array, first element is
// longitude and second is latitude
type Coordinates [2]float32

// BoundingBox is a set of two coordinates, the top left and the
// bottom right of a box
type BoundingBox [2]Coordinates

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates Coordinates `json:"coordinates"`
}

type CourseProperties struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	ZipCode string `json:"zipCode"`
}

type Course struct {
	Type       string           `json:"type"`
	Geometry   Geometry         `json:"geometry"`
	Properties CourseProperties `json:"properties"`
}

type CoursesDB struct {
	cluster *gocb.Cluster
	scope   *gocb.Scope
}

func NewCoursesDB(connectionString string, config *types.Config) *CoursesDB {
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: config.DbUser,
			Password: config.DbPassword,
		},
	}

	var connString string

	if connectionString != "" {
		connString = connectionString
	} else {
		connString = "127.0.0.1"
	}

	cluster, err := gocb.Connect(connString, options)

	scope := cluster.Bucket("courses").Scope("courses")

	if err != nil {
		log.Fatal(err)
	}

	if err := cluster.WaitUntilReady(5*time.Second, &gocb.WaitUntilReadyOptions{}); err != nil {
		log.Fatal(err)
	}

	return &CoursesDB{cluster, scope}
}

func (db *CoursesDB) GetCourseById(id string) (*Course, error) {
	query := "SELECT course.* from courses course WHERE course.properties.id = $id"

	params := make(map[string]interface{}, 1)
	params["id"] = id

	result, err := db.scope.Query(query, &gocb.QueryOptions{NamedParameters: params, Adhoc: true})

	if err != nil {
		log.Print(err)
		return &Course{}, err
	}

	course := &Course{}

	if err := result.One(&course); err != nil {
		return &Course{}, err
	}

	if (*course == Course{}) {
		return &Course{}, errors.New("not_found")
	}

	return course, nil
}

func (db *CoursesDB) ListCourses(limit int64) (*[]Course, error) {
	query := "SELECT course.* from courses course LIMIT $limit"

	params := make(map[string]interface{})
	params["limit"] = limit

	result, err := db.scope.Query(query, &gocb.QueryOptions{NamedParameters: params, Adhoc: true})

	if err != nil {
		log.Print(err)
		return nil, err
	}

	courses := make([]Course, 0)

	for result.Next() {
		var course *Course

		if err := result.Row(&course); err != nil {
			return nil, err
		}

		courses = append(courses, *course)
	}

	return &courses, nil
}

func (db *CoursesDB) SearchCoursesByBoundingBox(boundingBox BoundingBox) ([]Course, error) {
	mappedBoundingBox, err := mapBoundingBox(boundingBox)

	if err != nil {
		return nil, err
	}

	query := "select course.* from courses course WHERE course.geometry.coordinates[0] BETWEEN $topLng AND $bottomLng AND course.geometry.coordinates[1] BETWEEN $topLat AND $bottomLat"
	params := make(map[string]interface{}, 4)
	params["topLng"] = mappedBoundingBox[0][0]
	params["topLat"] = mappedBoundingBox[0][1]
	params["bottomLng"] = mappedBoundingBox[1][0]
	params["bottomLat"] = mappedBoundingBox[1][1]

	result, err := db.scope.Query(query, &gocb.QueryOptions{NamedParameters: params, Adhoc: true})

	if err != nil {
		return nil, err
	}

	courses := make([]Course, 0)

	for result.Next() {
		var course *Course

		if err := result.Row(&course); err != nil {
			return nil, err
		}

		courses = append(courses, *course)
	}

	return courses, nil
}

func (db *CoursesDB) SearchCoursesByPostCode(zipCode string) ([]Course, error) {
	query := "SELECT course.* from courses course WHERE course.properties.zipCode = $zipCode"

	params := make(map[string]interface{}, 1)
	params["zipCode"] = zipCode

	result, err := db.scope.Query(query, &gocb.QueryOptions{NamedParameters: params, Adhoc: true})

	if err != nil {
		return nil, err
	}

	courses := make([]Course, 0)

	for result.Next() {
		var course *Course

		if err := result.Row(&course); err != nil {
			return nil, err
		}

		courses = append(courses, *course)
	}

	return courses, nil
}

func mapBoundingBox(boundingBox BoundingBox) (BoundingBox, error) {
	firstCoord := boundingBox[0]
	secondCoord := boundingBox[1]

	if firstCoord[0] < secondCoord[0] && firstCoord[1] < secondCoord[1] {
		return boundingBox, nil
	} else if firstCoord[0] > secondCoord[0] && firstCoord[1] > secondCoord[1] {
		return BoundingBox{
			secondCoord,
			firstCoord,
		}, nil
	} else {
		return BoundingBox{}, errors.New("invalid_bounding_box")
	}
}

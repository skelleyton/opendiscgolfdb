package databases

import (
	"errors"
	"github.com/couchbase/gocb/v2"
	"log"
	"time"
)

type CoursesDB struct {
	cluster *gocb.Cluster
}

func NewCoursesDB(connectionString string) *CoursesDB {
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "placeholder",
			Password: "placeholder",
		},
	}
	cluster, err := gocb.Connect("127.0.0.1", options)
	
	if err != nil {
		log.Fatal(err)
	}
	
	if err := cluster.WaitUntilReady(5*time.Second, &gocb.WaitUntilReadyOptions{}); err != nil {
		log.Fatal(err)
	}
	
	return &CoursesDB{cluster}
}

func (db *CoursesDB) GetCourseById(id string) (*Course, error) {
	query := "SELECT course.* from courses.courses.courses course WHERE course.properties.id = $id"
	
	params := make(map[string]interface{}, 1)
	params["id"] = id
	
	result, err := db.cluster.Query(query, &gocb.QueryOptions{NamedParameters: params,  Adhoc: true})
	
	if err != nil {
		log.Print(err)
		return &Course{}, err
	}
	
	var courses []Course
	
	for result.Next() {
		var course *Course
		
		err := result.Row(&course)
		
		if err != nil {
			return &Course{}, err
		}
		
		courses = append(courses, *course)
	}
	
	if len(courses) == 0 {
		return &Course{}, errors.New("not_found")
	}
	
	return &courses[0], nil
}
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	dbByte, err := os.ReadFile("./db.json")

	if err != nil {
		log.Fatal(err)
	}

	var db *[]Course

	if err := json.Unmarshal(dbByte, &db); err != nil {
		log.Fatal(err)
	}

	database := DB{db}
	controller := Controller{Database: &database}

	r.GET("/course/:courseId", controller.GetCourse)
	r.GET("/courses", controller.ListCourses)
	r.GET("/courses/search", controller.Search)

	if err := r.Run(); err != nil {
		log.Fatal("Failed to initialize server")
	}
}

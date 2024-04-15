package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database := NewDB("")
	controller := Controller{Database: database}

	r.GET("/course/:courseId", controller.GetCourse)
	r.GET("/courses", controller.ListCourses)
	r.GET("/courses/search", controller.Search)

	if err := r.Run(); err != nil {
		log.Fatal("Failed to initialize server")
	}
}

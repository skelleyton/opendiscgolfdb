package main

import (
	"log"

	"opendgdb/controllers"
	"opendgdb/databases"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database := databases.NewCoursesDB("")
	controllers.NewCourseController(r, database)

	if err := r.Run(); err != nil {
		log.Fatal("Failed to initialize server")
	}
}

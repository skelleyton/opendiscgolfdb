package main

import (
	"log"

	"opendgdb/controllers"
	"opendgdb/databases"
	"opendgdb/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := utils.NewDotenvConfig("")

	database := databases.NewCoursesDB("", config)
	controllers.NewCourseController(r, database)

	if err := r.Run(); err != nil {
		log.Fatal("Failed to initialize server")
	}
}

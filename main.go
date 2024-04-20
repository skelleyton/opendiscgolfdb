package main

import (
	"log"

	"opendgdb/controllers"
	"opendgdb/databases"
	"opendgdb/types"
	"opendgdb/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := &types.Config{}

	utils.NewDotenvConfig("", config)

	database := databases.NewCoursesDB(config.ConnStr, config)
	controllers.NewCourseController(r, database)

	if err := r.Run(); err != nil {
		log.Fatal("Failed to initialize server")
	}
}

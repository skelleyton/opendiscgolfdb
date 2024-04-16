package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database := NewDB("")
	NewCourseController(r, database)

	if err := r.Run(); err != nil {
		log.Fatal("Failed to initialize server")
	}
}

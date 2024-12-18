package main

import (
	database "github.com/srad/docker-database-volume-backup/datbase"
	"github.com/srad/docker-database-volume-backup/services"
	"github.com/srad/docker-database-volume-backup/web"
	"os"
)

// @title API
// @version 1.0
// @description Backup server API
// @BasePath /api
func main() {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	database.Init()
	services.StartCron()
	web.ListenAndServe(httpPort)
}

package web

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/srad/docker-database-volume-backup/conf"
	routes2 "github.com/srad/docker-database-volume-backup/web/routes"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func ListenAndServe(port string) {
	config := conf.LoadConfig()
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Basic auth?
	useBasicAuth := config.BasicAuthUser != "" && config.BasicAuthPassword != ""
	if useBasicAuth {
		log.Infoln("Basic auth enabled")
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks.
			if subtle.ConstantTimeCompare([]byte(username), []byte(config.BasicAuthUser)) == 1 && subtle.ConstantTimeCompare([]byte(password), []byte(config.BasicAuthPassword)) == 1 {
				return true, nil
			}
			return false, nil
		}))
	} else {
		log.Infoln("Basic auth disabled")
	}

	e.Static("/backupfiles", "/backups")

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Index:  "index.html",
		HTML5:  true,
		Browse: false,
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler) // Serve Swagger UI

	api := e.Group("/api")

	api.GET("/backups", routes2.GetBackups)
	api.POST("/backups", routes2.CreateBackups)
	api.POST("/backups/:id/restore", routes2.RestoreBackup)
	api.GET("/backups/dumps", routes2.GetDumps)
	api.GET("/backups/files", routes2.GetFiles)
	api.DELETE("/backups/:id", routes2.DeleteBackup)

	api.GET("/cron", routes2.GetCron)

	api.GET("/config", func(c echo.Context) error {
		return c.JSON(http.StatusOK, config)
	})

	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}

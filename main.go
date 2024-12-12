package main

import (
    "crypto/subtle"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/robfig/cron/v3"
    "github.com/srad/wordpress-backup-enhanced/conf"
    "github.com/srad/wordpress-backup-enhanced/routes"
    "log"
    "net/http"
    "os"
)

var (
    appConfig    conf.AppConfig
    cronInstance *cron.Cron
)

func main() {
    appConfig = conf.LoadConfig()
    run()
}

func run() {
    cronInstance = cron.New()

    if appConfig.BackupOnStart {
        log.Println("Starting backup ...")
        go dump()
        go routes.ZipVolume()
    }

    log.Printf("Cron backup time defined: %s", appConfig.Cron)
    cronInstance.AddFunc(appConfig.Cron, func() { dump() })
    cronInstance.Start()
    listen()
}

func dump() {
    routes.Mysqldump(routes.MySqlConfig{Host: appConfig.Host, User: appConfig.User, Password: appConfig.Password, Database: appConfig.Database})
}

func listen() {
    e := echo.New()
    e.HideBanner = true

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    useBasicAuth := appConfig.BasicAuthUser != "" && appConfig.BasicAuthPassword != ""
    if useBasicAuth {
        log.Println("Basic auth enabled")
        e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
            // Be careful to use constant time comparison to prevent timing attacks.
            if subtle.ConstantTimeCompare([]byte(username), []byte(appConfig.BasicAuthUser)) == 1 && subtle.ConstantTimeCompare([]byte(password), []byte(appConfig.BasicAuthPassword)) == 1 {
                return true, nil
            }
            return false, nil
        }))
    } else {
        log.Println("Basic auth disabled")
    }

    e.Static("/backups", "/backups")

    e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
        Root:   "public",
        Index:  "index.html",
        HTML5:  true,
        Browse: false,
    }))

    api := e.Group("/api")

    api.GET("/dumps", routes.GetDumps)
    api.POST("/dumps", routes.CreateDump)

    api.GET("/files", routes.GetFiles)
    api.POST("/files", routes.CreateFiles)

    api.GET("/cron", func(c echo.Context) error {
        return routes.GetCron(c, cronInstance)
    })

    api.GET("/config", func(c echo.Context) error {
        return c.JSON(http.StatusOK, appConfig)
    })

    api.POST("/restore/:filename", func(c echo.Context) error {
        routes.MysqlRestoreDump(c.Param("filename"), routes.MySqlConfig{Host: appConfig.Host, User: appConfig.User, Password: appConfig.Password, Database: appConfig.Database})
        return c.JSON(http.StatusOK, nil)
    })

    httpPort := os.Getenv("PORT")
    if httpPort == "" {
        httpPort = "8080"
    }

    e.Logger.Fatal(e.Start("0.0.0.0:" + httpPort))
}

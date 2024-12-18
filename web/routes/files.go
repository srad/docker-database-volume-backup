package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/srad/docker-database-volume-backup/services"
	"net/http"
)

// GetFiles godoc
// @Summary Get list of backup files
// @Description Retrieves a list of all backup files from the backup directory.
// @Tags Backup Files
// @Accept json
// @Produce json
// @Success 200 {array} []string "List of backup files"
// @Failure 400 {object} HTTPError "Bad Request - Unable to retrieve backup files"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /files [get]
func GetFiles(c echo.Context) error {
	files, err := services.GetBackupFiles()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, files)
}

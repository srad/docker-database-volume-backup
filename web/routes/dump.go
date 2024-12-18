package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/srad/docker-database-volume-backup/conf"
	"github.com/srad/docker-database-volume-backup/utils"
	"net/http"
	"path"
)

// GetDumps godoc
// @Summary Get list of MySQL dump files
// @Description Retrieves a list of MySQL dump files from the backup directory.
// @Tags Dumps
// @Accept json
// @Produce json
// @Success 200 {array} []string "List of MySQL dump files"
// @Failure 400 {object} HTTPError "Bad Request - Unable to retrieve dump files"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /dumps [get]
func GetDumps(c echo.Context) error {
	dumpPath := path.Join(conf.GetBackupPath(), "dumps")
	files, err := utils.FilterBackup(dumpPath, "mysqldump")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, files)
}

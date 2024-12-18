package routes

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/srad/docker-database-volume-backup/conf"
	database "github.com/srad/docker-database-volume-backup/datbase"
	"github.com/srad/docker-database-volume-backup/services"
	"net/http"
	"strconv"
	"time"
)

type BackupDto struct {
	BackupId uint `json:"backupId"`

	DatabaseDumpFilePath string `json:"databaseDumpFilePath"`
	DatabaseDumpFilename string `json:"databaseDumpFilename"`
	DatabaseDumpFileSize uint64 `json:"databaseDumpFileSize"`

	VolumeFilePath     string `json:"volumeFilePath"`
	VolumeDumpFilename string `json:"volumeDumpFilename"`
	VolumeDumpFileSize uint64 `json:"volumeDumpFileSize"`

	CreatedAt time.Time `json:"createdAt"`
}

// GetBackups godoc
// @Summary Get list of backups
// @Description Retrieves a list of all available backups from the database.
// @Tags Backups
// @Accept json
// @Produce json
// @Success 200 {array} []BackupDto "List of backups"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /backups [get]
func GetBackups(c echo.Context) error {
	backups, err := database.BackupsList()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var backupList []BackupDto
	for _, backup := range backups {
		backupList = append(backupList, BackupDto{
			BackupId:             backup.BackupId,
			DatabaseDumpFilePath: backup.DatabaseDumpFilePath,
			DatabaseDumpFilename: backup.DatabaseDumpFilename,
			DatabaseDumpFileSize: backup.DatabaseDumpFileSize,
			VolumeFilePath:       backup.VolumeDumpFilePath,
			VolumeDumpFilename:   backup.VolumeDumpFilename,
			VolumeDumpFileSize:   backup.VolumeDumpFileSize,
			CreatedAt:            backup.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, backupList)
}

// CreateBackups godoc
// @Summary Create a new backup
// @Description Creates a new backup of the database using the provided configuration.
// @Tags Backups
// @Accept json
// @Produce json
// @Success 201 {string} string "Backup created"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /backups [post]
func CreateBackups(c echo.Context) error {
	config := conf.LoadConfig()
	mysqlConfig := config.ToMySqlConfig()

	if backup, err := services.CreateBackup(mysqlConfig); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		log.Infoln(backup.Info())
		return c.JSON(http.StatusCreated, "backup created")
	}
}

// RestoreBackup godoc
// @Summary Restore a backup
// @Description Restores a specific backup by its ID.
// @Tags Backups
// @Param id path uint32 true "Backup ID"
// @Accept json
// @Produce json
// @Success 200 {string} string "Backup restored successfully"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /backups/{id}/restore [post]
func RestoreBackup(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	config := conf.LoadConfig()
	mysqlConfig := config.ToMySqlConfig()

	if err := services.RestoreBackup(uint(id), mysqlConfig); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}

// DeleteBackup godoc
// @Summary Delete a backup
// @Description Deletes a specific backup by its ID.
// @Tags Backups
// @Param id path uint32 true "Backup ID"
// @Accept json
// @Produce json
// @Success 204 {string} string "Backup deleted successfully"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /backups/{id} [delete]
func DeleteBackup(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := services.DeleteBackup(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

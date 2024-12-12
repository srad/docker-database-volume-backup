package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/srad/wordpress-backup-enhanced/conf"
	"log"
	"net/http"
	"os/exec"
	"path"
	"time"
)

func CreateFiles(c echo.Context) error {
	ZipVolume()
	return c.JSON(http.StatusOK, nil)
}

func GetFiles(c echo.Context) error {
	filePath := path.Join(conf.GetBackupPath(), "files")
	files, err := filterBackup(filePath, "volume")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, files)
}

func ZipVolume() {
	sourceFolder := "/files"
	now := time.Now()
	stamp := now.Format("2006_02_01_15_04_05")
	commandString := fmt.Sprintf("tar czf \"/backups/files/volume_%s.tar.gz\" %s", stamp, sourceFolder)

	cmd := exec.Command("bash", "-c", commandString)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("Error running command: %s, %s", commandString, err.Error())
		return
	}

	log.Println("volumes backup completed")

	// Print the output
	log.Println(string(stdout))
}

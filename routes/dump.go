package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type MySqlConfig struct {
	Host     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Database string `validate:"required"`
}

type DumpInfo struct {
	Filename string    `json:"filename"`
	FileSize int64     `json:"fileSize"`
	Created  time.Time `json:"created"`
}

func GetDumps(c echo.Context) error {
	files, err := ListDumps()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, files)
}

func ListDumps() ([]DumpInfo, error) {
	backups := os.Getenv("BACKUP_FOLDER")
	if backups == "" {
		backups = "/backups"
	}

	entries, err := os.ReadDir(backups)
	if err != nil {
		log.Fatal(err)
	}

	var files []DumpInfo
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "mysqldump") {
			info, err := e.Info()
			if err != nil {
				return nil, err
			}
			files = append(files, DumpInfo{
				Filename: e.Name(),
				FileSize: info.Size(),
				Created:  info.ModTime(),
			})
		}
	}

	return files, nil
}

func Mysqldump(config MySqlConfig) {
	now := time.Now()
	stamp := now.Format("2006_02_01_15_04_05")
	commandString := fmt.Sprintf("mysqldump --host=\"%s\" --add-drop-table --no-tablespaces --user=\"%s\" --password=\"%s\" %s --single-transaction | bzip2 -c > \"/backups/mysqldump_%s.sql.bz2\"", config.Host, config.User, config.Password, config.Database, stamp)

	cmd := exec.Command("bash", "-c", commandString)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("Error running command: %s, %s", commandString, err.Error())
		return
	}

	log.Println("mysqldump completed")

	// Print the output
	log.Println(string(stdout))
}

func MysqlRestoreDump(file string, config MySqlConfig) {
	commandString := fmt.Sprintf("bunzip2 < \"%s\" | mysql --user=\"%s\" --password=\"%s\" --host=\"%s\" \"%s\"", file, config.User, config.Password, config.Host, config.Database)

	cmd := exec.Command("bash", "-c", commandString)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("Error running command: %s, %s", commandString, err.Error())
		return
	}

	log.Println("mysqldump completed")

	// Print the output
	log.Println(string(stdout))
}

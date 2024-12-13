package services

import (
	"fmt"
	"github.com/srad/docker-database-volume-backup/conf"
	"os/exec"
	"path"
)

func Mysqldump(config conf.MySqlConfig, stamp string) (string, string, error) {
	filename := fmt.Sprintf("mysqldump_%s.sql.bz2", stamp)
	filePath := path.Join("/backups/dumps", filename)
	commandString := fmt.Sprintf("mysqldump --host=\"%s\" --add-drop-table --no-tablespaces --user=\"%s\" --password=\"%s\" %s --single-transaction | bzip2 -c > \"%s\"", config.Host, config.User, config.Password, config.Database, filePath)

	cmd := exec.Command("bash", "-c", commandString)
	stdout, err := cmd.Output()

	if err != nil {
		return "", "", fmt.Errorf("error running command '%s': %s, %s", commandString, err.Error(), string(stdout))
	}

	return filePath, filename, nil
}

func MysqlRestoreDump(filePath string, config conf.MySqlConfig) error {
	commandString := fmt.Sprintf("bunzip2 < \"%s\" | mysql --user=\"%s\" --password=\"%s\" --host=\"%s\" \"%s\"", filePath, config.User, config.Password, config.Host, config.Database)

	cmd := exec.Command("bash", "-c", commandString)
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("error running command '%s': %s, %s", commandString, err.Error(), string(stdout))
	}

	return nil
}

package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path"
)

func CheckTarArchive(filePath string) (string, error) {
	commandString := fmt.Sprintf("tar -tzf %s >/dev/null", filePath)

	cmd := exec.Command("bash", "-c", commandString)
	output, err := cmd.CombinedOutput()
	str := string(output)

	if err != nil {
		return "", fmt.Errorf("error running command '%s': %s, %s", commandString, err.Error(), str)
	}

	return str, nil
}

func ZipPath(folderToZip, folderToStoreZip, stamp string) (string, string, error) {
	filename := fmt.Sprintf("volume_%s.tar.gz", stamp)
	zipFolderFilePath := path.Join(folderToStoreZip, filename)

	commandString := fmt.Sprintf("tar czf %s %s", zipFolderFilePath, folderToZip)

	log.Infof("Executing command: %s", commandString)

	cmd := exec.Command("bash", "-c", commandString)
	output, err := cmd.CombinedOutput()
	str := string(output)

	if err != nil {
		return "", "", fmt.Errorf("error running command '%s': %s, %s", commandString, err.Error(), str)
	}

	return zipFolderFilePath, filename, nil
}

func UnzipFile(zipFilePath, unzipFolder string) error {
	commandString := fmt.Sprintf("tar -xzf \"%s\" --directory=\"%s\"", zipFilePath, unzipFolder)

	cmd := exec.Command("bash", "-c", commandString)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

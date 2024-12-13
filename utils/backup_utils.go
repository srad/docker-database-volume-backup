package utils

import (
	"os"
	"strings"
	"time"
)

type BackupInfo struct {
	Filename string    `json:"filename"`
	FileSize int64     `json:"fileSize"`
	Created  time.Time `json:"created"`
}

func FilterBackup(folder, filePrefix string) ([]BackupInfo, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var files []BackupInfo
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "mysqldump") {
			info, err := e.Info()
			if err != nil {
				return nil, err
			}
			files = append(files, BackupInfo{
				Filename: e.Name(),
				FileSize: info.Size(),
				Created:  info.ModTime(),
			})
		}
	}

	return files, nil
}

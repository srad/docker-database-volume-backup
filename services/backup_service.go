package services

import (
	"errors"
	"github.com/srad/docker-database-volume-backup/conf"
	database "github.com/srad/docker-database-volume-backup/datbase"
	"github.com/srad/docker-database-volume-backup/utils"
	"os"
	"path"
	"time"
)

type BackupResult struct {
	MysqlDumpFilePath string `json:"mysqlDumpFilePath"`
	MysqlDumpFilename string `json:"mysqlDumpFilename"`

	VolumeDumpFilepath string `json:"volumeDumpFilePath"`
	VolumeDumpFilename string `json:"volumeDumpFilename"`
}

func GetBackupFiles() ([]utils.BackupInfo, error) {
	filePath := path.Join(conf.GetBackupPath(), "files")
	files, err := utils.FilterBackup(filePath, "volume")

	if err != nil {
		return nil, err
	}

	return files, err
}

func CreateBackup(mysqlConf conf.MySqlConfig) (*database.Backup, error) {
	now := time.Now()

	result, err := dumpFiles(mysqlConf, now)
	if err != nil {
		return nil, err
	}

	dbDumpSize := uint64(0)
	volumeDumpSize := uint64(0)

	if fileInfo, err := os.Stat(result.VolumeDumpFilepath); err == nil {
		volumeDumpSize = uint64(fileInfo.Size())
	}

	if fileInfo, err := os.Stat(result.MysqlDumpFilePath); err == nil {
		dbDumpSize = uint64(fileInfo.Size())
	}

	backup := &database.Backup{
		DatabaseDumpFilePath: result.MysqlDumpFilePath,
		DatabaseDumpFilename: result.MysqlDumpFilename,
		DatabaseDumpFileSize: dbDumpSize,

		VolumeDumpFilePath: result.VolumeDumpFilepath,
		VolumeDumpFilename: result.VolumeDumpFilename,
		VolumeDumpFileSize: volumeDumpSize,

		CreatedAt: now,
	}

	return backup, backup.SaveBackup()
}

func RestoreBackup(backupId uint, mysqlConf conf.MySqlConfig) error {
	backup, errFind := database.FindBackupById(backupId)
	if errFind != nil {
		return errFind
	}

	// 1. Delete folders in volume mapping.
	if err := os.RemoveAll(conf.GetVolumeMapping()); err != nil {
		return err
	}

	// 2. Restore database
	if err := MysqlRestoreDump(backup.DatabaseDumpFilePath, mysqlConf); err != nil {
		return err
	}

	// 3. Restore files
	if err := utils.UnzipFile(backup.VolumeDumpFilePath, conf.GetBackupFilesPath()); err != nil {
		return err
	}

	return nil
}

func dumpFiles(mysqlConf conf.MySqlConfig, timestamp time.Time) (*BackupResult, error) {
	stamp := timestamp.Format("2006_02_01_15_04_05")

	zipFilePath, volumeDumpFilename, err1 := utils.ZipPath("/files", "/backups/files", stamp)
	dumpFilePath, dbDumpFilename, err2 := Mysqldump(mysqlConf, stamp)

	// If one fails then cancel both ...
	if err1 != nil || err2 != nil {
		_ = os.Remove(dumpFilePath)
		_ = os.Remove(zipFilePath)
		return nil, errors.Join(err1, err2)
	}

	result := &BackupResult{
		MysqlDumpFilePath:  dumpFilePath,
		MysqlDumpFilename:  dbDumpFilename,
		VolumeDumpFilepath: zipFilePath,
		VolumeDumpFilename: volumeDumpFilename,
	}

	return result, nil
}

func DeleteBackup(id uint) error {
	backup, errFind := database.FindBackupById(id)
	if errFind != nil {
		return errFind
	}

	if err := os.RemoveAll(backup.VolumeDumpFilePath); err != nil {
		return err
	}
	if err := os.RemoveAll(backup.DatabaseDumpFilePath); err != nil {
		return err
	}
	if err := database.DeleteBackupById(id); err != nil {
		return err
	}
	return nil
}

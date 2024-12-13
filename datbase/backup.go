package database

import (
	"fmt"
	"time"
)

type Backup struct {
	BackupId uint `gorm:"autoIncrement;primaryKey;column:backup_id"`

	DatabaseDumpFilePath string `gorm:"not null;default:null;column:database_dump_file_path"`
	DatabaseDumpFileSize uint64 `gorm:"column:database_dump_file_size"`
	DatabaseDumpFilename string `gorm:"not null;default:null;column:database_dump_filename"`

	VolumeDumpFilePath string `gorm:"not null;default:null;column:volume_dump_file_path"`
	VolumeDumpFileSize uint64 `gorm:"column:volume_dump_file_size"`
	VolumeDumpFilename string `gorm:"not null;default:null;column:volume_dump_filename"`

	CreatedAt time.Time `gorm:"not null;default:null;column:created_at"`
}

type BackupFile struct {
	Filename string
}

func (backup *Backup) SaveBackup() error {
	return dbInstance.Model(Backup{}).Save(backup).Error
}

func BackupsList() ([]Backup, error) {
	var backups []Backup
	err := dbInstance.Find(&backups).Error
	if err != nil {
		return nil, err
	}
	return backups, nil
}

func FindBackupById(id uint) (*Backup, error) {
	var backup Backup
	err := dbInstance.First(&backup, id).Error
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

func DeleteBackupById(id uint) error {
	return dbInstance.Delete(&Backup{}, id).Error
}

func (backup *Backup) Info() string {
	return fmt.Sprintf("Database: %s (%d MB), Volume: %s (%d MB)", backup.DatabaseDumpFilePath, backup.DatabaseDumpFileSize/1024/1024, backup.VolumeDumpFilePath, backup.VolumeDumpFileSize/1024/1024)
}

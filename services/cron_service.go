package services

import (
	"errors"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/srad/docker-database-volume-backup/conf"
	"time"
)

type CronEntry struct {
	Id   cron.EntryID
	Next time.Time
}

var cronInstance *cron.Cron

func StartCron() {
	config := conf.LoadConfig()

	if config.BackupOnStart {
		log.Println("Starting backup ...")
		config := conf.LoadConfig()
		if backup, err := CreateBackup(config.ToMySqlConfig()); err != nil {
			log.Fatal(err)
		} else {
			log.Infoln("Backup created successfully!")
            log.Infoln(backup.Info())
		}
	}

	if cronInstance == nil {
		cronInstance = cron.New()
	}

	log.Printf("Starting cron process with spec: %s", config.Cron)
	cronInstance.AddFunc(config.Cron, func() { CreateBackup(config.ToMySqlConfig()) })
	cronInstance.Start()
}

func StopCron() {
	if cronInstance != nil {
		cronInstance.Stop()
	}
}

func GetCronJobs() ([]CronEntry, error) {
	if cronInstance == nil {
		return nil, errors.New("cron instance not initialized")
	}

	entries := cronInstance.Entries()
	var result = make([]CronEntry, len(entries))
	for _, entry := range entries {
		result = append(result, CronEntry{
			Id:   entry.ID,
			Next: entry.Next,
		})
	}

	return result, nil
}

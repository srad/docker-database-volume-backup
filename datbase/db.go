package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/srad/docker-database-volume-backup/conf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB

func Init() {
	config := conf.LoadConfig()

	log.Infof("Creating database at location: %s", config.DatabaseFilePath)

	newLogger := logger.New(
		log.New(),
		logger.Config{
			//SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			//ParameterizedQueries:      true,         // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)

	dialector := sqlite.Open(config.DatabaseFilePath)

	/// Open and assign database.
	dbConfig := &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	db, err := gorm.Open(dialector, dbConfig)
	if err != nil {
		panic("failed to connect models")
	}
	dbInstance = db

	migrate()
}

func migrate() {
	// Migrate the schema
	if err := dbInstance.AutoMigrate(&Backup{}); err != nil {
		panic(fmt.Sprintf("[Migrate] Error user: %s", err))
	}
}

package core

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase returns a gorm.DB struct, gorm.DB.DB() returns a database handle
// see http://golang.org/pkg/database/sql/#DB
func NewDatabase(cfg *Config, dbName string) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger: gormLogger,
	}

	db, err := gorm.Open(
		sqlite.Open(":memory:"),
		gormConfig,
	)
	if cfg.Debug {
		db = db.Debug()
	}
	if err != nil {
		return db, err
	}

	// if _, err := db.DB(); err != nil {
	// 	return db, err
	// }

	// sqlDb.Set("gorm:table_options", "charset=ascii")
	// // Database logging
	// sqlDb.LogMode(cfg.Debug)

	return db, nil
}

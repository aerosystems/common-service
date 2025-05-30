package gormclient

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	gormv2logrus "github.com/thomas-tacquet/gormv2-logrus"
	"gorm.io/driver/postgres"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func NewPostgresDB(log *logrus.Logger, cfg *Config) *gorm.DB {
	once.Do(func() {
		e := logrus.NewEntry(log)

		gormLogger := gormv2logrus.NewGormlog(gormv2logrus.WithLogrusEntry(e))
		gormLogger.LogMode(logger.Info)
		gormLogger.SlowThreshold = 100 * time.Millisecond
		gormLogger.SkipErrRecordNotFound = true
		count := 0
		for {
			db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
				Logger:                 gormLogger,
				SkipDefaultTransaction: true,
				PrepareStmt:            true,
			})
			if err != nil {
				e.Logger.Info("PostgreSQL not ready...")
				count++
			} else {
				e.Logger.Info("Connected to PostgreSQL.")
				instance = db
				return
			}
			if count > 10 {
				e.Logger.Info(err)
				return
			}
			e.Logger.Info("Backing off for two seconds...")
			time.Sleep(2 * time.Second)
		}
	})
	return instance
}

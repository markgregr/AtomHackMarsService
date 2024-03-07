package db

import (
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct{
	DatabaseGORM *gorm.DB
	DatabaseCfg *config.Database
}

func (d *Database) New(cfg *config.Database) error {
	var err error
	DatabaseCfg := cfg
	l := logger.Default
	if log.StandardLogger().GetLevel() >= log.DebugLevel {
		l = l.LogMode(logger.Info)
	}
	log.Println(DatabaseCfg.DSN)
	d.DatabaseGORM, err = gorm.Open(postgres.Open(DatabaseCfg.DSN), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		log.Info("execute database connection")
		return err
	}

	if cfg.AutoMigrate {
		log.Info("execute database migrations")
		if err := Migrate(d.DatabaseGORM); err != nil {
			return err
		}
	}

	return nil
}
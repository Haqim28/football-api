package database

import (
	"log"

	"github.com/yourname/football-api/internal/config"
	"github.com/yourname/football-api/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	log.Println("database connected and migrated")
	return db
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Team{},
		&domain.Player{},
		&domain.Match{},
		&domain.Goal{},
	)
}

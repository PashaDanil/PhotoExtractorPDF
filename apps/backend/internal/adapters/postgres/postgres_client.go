package postgres

import (
	"api/pkg/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.DBName,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.SSLMode,
		cfg.PostgresConfig.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close postgres connection: %w", err)
	}

	return nil
}

func (p *Postgres) GetDB() *gorm.DB {
	return p.db
}

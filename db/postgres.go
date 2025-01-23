package db

import (
	"context"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DNSLog struct {
	ID        uint   `gorm:"primaryKey"`
	Timestamp int64  `gorm:"autoCreateTime"`
	Domain    string `gorm:"index"`
	QueryType string
	Response  string
}

type PostgresDB struct {
	Client *gorm.DB
}

func NewPostgresDB(dsn string) *PostgresDB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	err = db.AutoMigrate(&DNSLog{})
	if err != nil {
		log.Fatalf("Failed to migrate database : %v", err)
	}

	return &PostgresDB{Client: db}
}

func (p *PostgresDB) LogQuery(ctx context.Context, domain, queryType, response string) error {
	log := DNSLog{
		Domain:    domain,
		QueryType: queryType,
		Response:  response,
	}

	return p.Client.WithContext(ctx).Create(&log).Error
}

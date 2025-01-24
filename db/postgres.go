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

type DNSRecord struct {
	ID     uint   `gorm:"primaryKey"`
	Domain string `gorm:"index;not null"`
	Type   string `gorm:"not null"`
	Value  string `gorm:"not null"`
	TTL    int    `gorm:"not null"`
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
		log.Fatalf("Failed to migrate database(Log) : %v", err)
	}

	err = db.AutoMigrate(&DNSRecord{})
	if err != nil {
		log.Fatalf("Failed to migrate database(Record) : %v", err)
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

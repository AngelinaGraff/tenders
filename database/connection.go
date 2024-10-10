package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func ConnectDatabase(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
        return nil, err
    }
    return db, nil
}

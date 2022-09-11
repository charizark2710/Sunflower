package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func DbConfig() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=" + os.Getenv("POSTGRES_HOST") +
			" user=" + os.Getenv("POSTGRES_USER") +
			" password=" + os.Getenv("POSTGRES_PASSWORD") +
			" dbname=" + os.Getenv("POSTGRES_DB") +
			" port=" + os.Getenv("POSTGRES_PORT") +
			" sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: os.Getenv("POSTGRES_SCHEMA") + ".", // schema name
		}})
	if err != nil {
		return nil, err
	}
	err = db.Exec("CREATE SCHEMA IF NOT EXISTS " + os.Getenv("POSTGRES_SCHEMA")).Error
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error

	return db, err
}

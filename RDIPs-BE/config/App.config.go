package config

import (
	"database/sql"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	AMQP_handler "RDIPs-BE/handler/AMQP"
)

func DbConfig() (*gorm.DB, error) {
	sql := &sql.DB{}
	sql.SetMaxIdleConns(10)
	sql.SetConnMaxLifetime(1 * time.Minute)
	sql.SetConnMaxIdleTime(time.Second * 30)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=" + os.Getenv("POSTGRES_HOST") +
			" user=" + os.Getenv("POSTGRES_USER") +
			" password=" + os.Getenv("POSTGRES_PASSWORD") +
			" dbname=" + os.Getenv("POSTGRES_DB") +
			" port=" + os.Getenv("POSTGRES_PORT") +
			" sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
	}), &gorm.Config{
		ConnPool: sql,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: os.Getenv("POSTGRES_SCHEMA") + ".", // schema name
		}})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	err = db.Exec("CREATE SCHEMA IF NOT EXISTS " + os.Getenv("POSTGRES_SCHEMA")).Error
	return db, err
}

func RabbitMqConfig() error {
	return AMQP_handler.InitializeAMQP()
}

package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	AMQPconst "RDIPs-BE/constant/AMQP_Const"

	amqp "github.com/rabbitmq/amqp091-go"
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

	return db, err
}

func RabbitMqConfig() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(
		"amqp://" +
			os.Getenv("BROKER_USER") +
			":" + os.Getenv("BROKER_PASSWORD") +
			"@" + os.Getenv("BROKER_HOST") +
			":" + os.Getenv("BROKER_PORT") + "/")
	if err != nil {
		return nil, nil, err
	}

	ch, chErr := conn.Channel()
	if chErr != nil {
		return nil, nil, chErr
	}

	// Declare exchange
	for _, exchange := range AMQPconst.ExhangeArr {
		ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	}

	return conn, ch, nil
}

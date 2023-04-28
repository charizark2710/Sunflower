package model

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type helper struct {
	db   *gorm.DB
	amqp AmqpModel
}

var Helper *helper = &helper{}

func (h *helper) GetDb() *gorm.DB {
	return h.db
}

func (h *helper) SetDb(db *gorm.DB) {
	h.db = db
}

func (h *helper) GetAMQPConnection() *amqp.Connection {
	return h.amqp.conn
}

func (h *helper) SetAMQPChannel(ch *amqp.Channel) {
	h.amqp.channel = ch
}

func (h *helper) GetAMQPChannel() *amqp.Channel {
	return h.amqp.channel
}

func (h *helper) SetAMQP(conn *amqp.Connection) {
	h.amqp.conn = conn
}

package model

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpModel struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

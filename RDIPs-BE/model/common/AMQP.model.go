package model

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BaseAmqpConn interface {
	Channel() (*amqp.Channel, error)
	Close() error
	CloseDeadline(deadline time.Time) error
	ConnectionState() tls.ConnectionState
	IsClosed() bool
	LocalAddr() net.Addr
	NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
	RemoteAddr() net.Addr
	UpdateSecret(newSecret string, reason string) error
}

type BaseAmqpChannel interface {
	Ack(tag uint64, multiple bool) error
	Cancel(consumer string, noWait bool) error
	Close() error
	Confirm(noWait bool) error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeBind(destination string, key string, source string, noWait bool, args amqp.Table) error
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error
	ExchangeDeclarePassive(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error
	ExchangeDelete(name string, ifUnused bool, noWait bool) error
	ExchangeUnbind(destination string, key string, source string, noWait bool, args amqp.Table) error
	Flow(active bool) error
	Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error)
	GetNextPublishSeqNo() uint64
	IsClosed() bool
	Nack(tag uint64, multiple bool, requeue bool) error
	NotifyCancel(c chan string) chan string
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	NotifyConfirm(ack chan uint64, nack chan uint64) (chan uint64, chan uint64)
	NotifyFlow(c chan bool) chan bool
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	NotifyReturn(c chan amqp.Return) chan amqp.Return
	Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
	PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
	PublishWithDeferredConfirm(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error)
	PublishWithDeferredConfirmWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error)
	Qos(prefetchCount int, prefetchSize int, global bool) error
	QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDelete(name string, ifUnused bool, ifEmpty bool, noWait bool) (int, error)
	QueueInspect(name string) (amqp.Queue, error)
	QueuePurge(name string, noWait bool) (int, error)
	QueueUnbind(name string, key string, exchange string, args amqp.Table) error
	Recover(requeue bool) error
	Reject(tag uint64, requeue bool) error
	Tx() error
	TxCommit() error
	TxRollback() error
}

type AmqpModel struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

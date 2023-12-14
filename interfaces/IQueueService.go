package interfaces

import amqp "github.com/rabbitmq/amqp091-go"

type IQueueService interface {
	PushDataToQueue(data []byte, queueName string) error
	ConsumeQueue(queueName string) <-chan amqp.Delivery
}

package infrastructures

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func NewQueueService(connectionString string, queueNames []string) *QueueService {
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	service := &QueueService{
		connection: conn,
		ch:         ch,
	}

	service.InitQueues(queueNames)

	return service
}

type QueueService struct {
	connection *amqp.Connection
	ch         *amqp.Channel
	queues     []amqp.Queue
}

func (s *QueueService) InitQueues(queueNames []string) {
	for _, name := range queueNames {
		q, err := s.ch.QueueDeclare(
			name,  // name
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		failOnError(err, "Failed to declare a queue "+name)

		s.queues = append(s.queues, q)
	}
}

func (s *QueueService) PushDataToQueue(data []byte, queueName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        data,
		})

	return err
}

func (s *QueueService) ConsumeQueue(queueName string) <-chan amqp.Delivery {
	messages, err := s.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	return messages
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

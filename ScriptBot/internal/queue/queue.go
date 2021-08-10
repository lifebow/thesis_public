package queue

import (
	"fmt"
	"scriptbot/internal/config"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

const (
	NormalChannelName    = "normalCheck"
	EmergencyChannelName = "emergencyCheck"
)

type QueueConnection struct {
	JobConnection *amqp.Connection
	JobQueue      *amqp.Queue
	JobChannel    *amqp.Channel
	EmergencyQueue      *amqp.Queue
	EmergencyChannel *amqp.Channel
	mutex                 sync.Mutex
}

func (queueConnection *QueueConnection) Close() {
	queueConnection.JobChannel.Close()
	queueConnection.JobConnection.Close()
	log.Info().Msg("RabbitMQ disconnected")
}

func (queueConnection *QueueConnection) ConsumeJob() (<-chan amqp.Delivery, error) {
	queueConnection.mutex.Lock()
	defer queueConnection.mutex.Unlock()
	return queueConnection.JobChannel.Consume(
		queueConnection.JobQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}
func (queueConnection *QueueConnection) ConsumeEmergencyJob() (<-chan amqp.Delivery, error) {
	queueConnection.mutex.Lock()
	defer queueConnection.mutex.Unlock()
	return queueConnection.JobChannel.Consume(
		queueConnection.EmergencyQueue.Name,
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

func NewQueue(cfg config.WorkerConfig) (*QueueConnection, error) {
	username := cfg.RabbitmqDefaultUser
	password := cfg.RabbitmqDefaultPass
	rabbitmqHost := cfg.RabbitmqHost
	connectionUri := fmt.Sprintf("amqp://%v:%v@%v:5672/", username, password, rabbitmqHost)

	fmt.Println("RabbitmqUrl", connectionUri)

	timesTry := 0
	jobConnection, err := amqp.Dial(connectionUri)
	for err != nil && timesTry < 10 {
		time.Sleep(time.Second * 5)
		log.Info().Msg("Failed to connect to RabbitMQ, retrying ...")
		timesTry++
		jobConnection, err = amqp.Dial(connectionUri)
	}

	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	jobChannel, err := jobConnection.Channel()
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	q, err := jobChannel.QueueDeclare(
		NormalChannelName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	// get 4 jobs at a time
	err = jobChannel.Qos(
		4,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	log.Info().Msg("RabbitMQ connected")
	emergencyChannel, err := jobConnection.Channel()
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	err = emergencyChannel.ExchangeDeclare(
		EmergencyChannelName,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	q2, err := emergencyChannel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	err = emergencyChannel.QueueBind(
		q2.Name, // queue name
		"",     // routing key
		EmergencyChannelName, // exchange
		false,
		nil,
	)
	return &QueueConnection{
		JobConnection: jobConnection,
		JobQueue:      &q,
		JobChannel:    jobChannel,
		EmergencyQueue: &q2,
		EmergencyChannel: emergencyChannel,
	}, nil
}

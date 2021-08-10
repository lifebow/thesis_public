package queue

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type QueueConnection struct {
	Connection       *amqp.Connection
	NormalQueue      *amqp.Queue
	EmergencyQueue   *amqp.Queue
	NormalChannel    *amqp.Channel
	EmergencyChannel *amqp.Channel
	mutex            sync.Mutex
}

const (
	NormalChannelName    = "normalCheck"
	EmergencyChannelName = "emergencyCheck"
)

func (q *QueueConnection) Close() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	_ = q.NormalChannel.Close()
	_ = q.EmergencyChannel.Close()
	_ = q.Connection.Close()
	log.Info().Msg("RabbitMQ disconnected")
}

func (q *QueueConnection) PublishNormalJob(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	tmp := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         data,
	}
	err = q.NormalChannel.Publish(
		"",
		q.NormalQueue.Name,
		false,
		false,
		tmp,
	)
	if err != nil {
		log.Warn().Msg("Failed to publish new job, trying to reconnect to RabbitMQ ...")
		q.Close()
		times := 0
		for times < 3 {
			err = q.Connect()
			if err == nil {
				break
			}
			times++
		}
	}

	return err
}
func (q *QueueConnection) PublishEmergencyJob(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = q.EmergencyChannel.Publish(
		EmergencyChannelName,
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		},
	)
	if err != nil {
		log.Warn().Msg("Failed to publish new job, trying to reconnect to RabbitMQ ...")
		q.Close()
		times := 0
		for times < 3 {
			err = q.Connect()
			if err == nil {
				break
			}
			times++
		}
	}

	return err
}
func (q *QueueConnection) Connect() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	username := "guest"
	password := "guest"

	rabbitHost := "127.0.0.1"
	connectionUri := fmt.Sprintf("amqp://%v:%v@%v:5672/", username, password, rabbitHost)
	log.Info().Msgf(connectionUri)
	var err error
	times := 0

	for times < 10 {
		q.Connection, err = amqp.Dial(connectionUri)
		if err != nil {
			log.Info().Msg("Failed to connect to RabbitMQ, retrying ...")
			time.Sleep(time.Second * 5)
			times++
		} else {
			break
		}
	}
	if err != nil {
		log.Fatal().Msg("Failed to connect to RabbitMQ, please check again!")
		return err
	}
	q.NormalChannel, err = q.Connection.Channel()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	q1, err := q.NormalChannel.QueueDeclare(
		NormalChannelName,
		true,
		false,
		false,
		false,
		nil,
	)
	q.NormalQueue = &q1
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	q.EmergencyChannel, err = q.Connection.Channel()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	err = q.EmergencyChannel.ExchangeDeclare(
		EmergencyChannelName, // name
		"fanout",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	log.Info().Msg("RabbitMQ connected")
	return nil
}

func NewQueue() (*QueueConnection, error) {
	var queueConnection QueueConnection
	err := queueConnection.Connect()
	return &queueConnection, err
}

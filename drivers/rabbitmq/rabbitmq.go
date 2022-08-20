package rabbitmq

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type RabbitMQ struct {
	Client   *amqp.Connection
	URL      string
	Exchange string
	Queue    string
}

func (d *RabbitMQ) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "rabbitmq",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"RABBITMQ_URL") != "" {
		d.URL = os.Getenv(prefix + "RABBITMQ_URL")
	}
	if os.Getenv(prefix+"RABBITMQ_QUEUE") != "" {
		d.Queue = os.Getenv(prefix + "RABBITMQ_QUEUE")
	}
	return nil
}

func (d *RabbitMQ) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "rabbitmq",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.URL = *flags.RabbitMQURL
	d.Queue = *flags.RabbitMQQueue
	return nil
}

func (d *RabbitMQ) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "rabbitmq",
		"fn":  "Init",
	})
	l.Debug("Initializing rabbitmq driver")
	conn, err := amqp.Dial(d.URL)
	if err != nil {
		return err
	}
	d.Client = conn
	return nil
}

func (d *RabbitMQ) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "rabbitmq",
		"fn":  "Push",
	})
	l.Debug("Pushing message to rabbitmq")
	ch, err := d.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		d.Queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        bd,
	}
	err = ch.PublishWithContext(
		context.Background(),
		d.Exchange, // exchange
		q.Name,     // routing key
		false,      // mandatory
		false,      // immediate
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *RabbitMQ) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "rabbitmq",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up rabbitmq driver")
	if err := d.Client.Close(); err != nil {
		return err
	}
	return nil
}

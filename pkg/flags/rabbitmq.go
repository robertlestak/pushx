package flags

var (
	RabbitMQURL      = FlagSet.String("rabbitmq-url", "", "RabbitMQ URL")
	RabbitMQQueue    = FlagSet.String("rabbitmq-queue", "", "RabbitMQ queue")
	RabbitMQExchange = FlagSet.String("rabbitmq-exchange", "", "RabbitMQ exchange")
)

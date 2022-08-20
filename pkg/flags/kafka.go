package flags

var (
	KafkaBrokers      = FlagSet.String("kafka-brokers", "", "Kafka brokers, comma separated")
	KafkaTopic        = FlagSet.String("kafka-topic", "", "Kafka topic")
	KafkaEnableTLS    = FlagSet.Bool("kafka-enable-tls", false, "Enable TLS")
	KafkaTLSInsecure  = FlagSet.Bool("kafka-tls-insecure", false, "Enable TLS insecure")
	KafkaCAFile       = FlagSet.String("kafka-tls-ca-file", "", "Kafka TLS CA file")
	KafkaCertFile     = FlagSet.String("kafka-tls-cert-file", "", "Kafka TLS cert file")
	KafkaKeyFile      = FlagSet.String("kafka-tls-key-file", "", "Kafka TLS key file")
	KafkaEnableSasl   = FlagSet.Bool("kafka-enable-sasl", false, "Enable SASL")
	KafkaSaslType     = FlagSet.String("kafka-sasl-type", "", "Kafka SASL type. Can be either 'scram' or 'plain'")
	KafkaSaslUsername = FlagSet.String("kafka-sasl-username", "", "Kafka SASL user")
	KafkaSaslPassword = FlagSet.String("kafka-sasl-password", "", "Kafka SASL password")
)

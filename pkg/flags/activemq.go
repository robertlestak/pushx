package flags

var (
	ActiveMQAddress     = FlagSet.String("activemq-address", "", "ActiveMQ STOMP address")
	ActiveMQName        = FlagSet.String("activemq-name", "", "ActiveMQ name")
	ActiveMQEnableTLS   = FlagSet.Bool("activemq-enable-tls", false, "Enable TLS")
	ActiveMQTLSInsecure = FlagSet.Bool("activemq-tls-insecure", false, "Enable TLS insecure")
	ActiveMQTLSCA       = FlagSet.String("activemq-tls-ca-file", "", "TLS CA")
	ActiveMQTLSCert     = FlagSet.String("activemq-tls-cert-file", "", "TLS cert")
	ActiveMQTLSKey      = FlagSet.String("activemq-tls-key-file", "", "TLS key")
)

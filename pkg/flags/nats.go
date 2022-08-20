package flags

var (
	NATSURL         = FlagSet.String("nats-url", "", "NATS URL")
	NATSSubject     = FlagSet.String("nats-subject", "", "NATS subject")
	NATSCredsFile   = FlagSet.String("nats-creds-file", "", "NATS creds file")
	NATSJWTFile     = FlagSet.String("nats-jwt-file", "", "NATS JWT file")
	NATSNKeyFile    = FlagSet.String("nats-nkey-file", "", "NATS NKey file")
	NATSUsername    = FlagSet.String("nats-username", "", "NATS username")
	NATSPassword    = FlagSet.String("nats-password", "", "NATS password")
	NATSToken       = FlagSet.String("nats-token", "", "NATS token")
	NATSEnableTLS   = FlagSet.Bool("nats-enable-tls", false, "NATS enable TLS")
	NATSTLSInsecure = FlagSet.Bool("nats-tls-insecure", false, "NATS TLS insecure")
	NATSTLSCAFile   = FlagSet.String("nats-tls-ca-file", "", "NATS TLS CA file")
	NATSTLSCertFile = FlagSet.String("nats-tls-cert-file", "", "NATS TLS cert file")
	NATSTLSKeyFile  = FlagSet.String("nats-tls-key-file", "", "NATS TLS key file")
)

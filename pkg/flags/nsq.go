package flags

var (
	NSQNSQLookupdAddress = FlagSet.String("nsq-nsqlookupd-address", "", "NSQ nsqlookupd address")
	NSQNSQDAddress       = FlagSet.String("nsq-nsqd-address", "", "NSQ nsqd address")
	NSQTopic             = FlagSet.String("nsq-topic", "", "NSQ topic")
	NSQEnableTLS         = FlagSet.Bool("nsq-enable-tls", false, "Enable TLS")
	NSQTLSSkipVerify     = FlagSet.Bool("nsq-tls-skip-verify", false, "NSQ TLS skip verify")
	NSQCAFile            = FlagSet.String("nsq-tls-ca-file", "", "NSQ TLS CA file")
	NSQCertFile          = FlagSet.String("nsq-tls-cert-file", "", "NSQ TLS cert file")
	NSQKeyFile           = FlagSet.String("nsq-tls-key-file", "", "NSQ TLS key file")
)

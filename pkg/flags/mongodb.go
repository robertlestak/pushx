package flags

var (
	MongoHost        = FlagSet.String("mongo-host", "", "MongoDB host")
	MongoPort        = FlagSet.String("mongo-port", "27017", "MongoDB port")
	MongoUser        = FlagSet.String("mongo-user", "", "MongoDB user")
	MongoPassword    = FlagSet.String("mongo-password", "", "MongoDB password")
	MongoDatabase    = FlagSet.String("mongo-database", "", "MongoDB database")
	MongoCollection  = FlagSet.String("mongo-collection", "", "MongoDB collection")
	MongoEnableTLS   = FlagSet.Bool("mongo-enable-tls", false, "Enable TLS")
	MongoAuthSource  = FlagSet.String("mongo-auth-source", "", "MongoDB auth source")
	MongoTLSInsecure = FlagSet.Bool("mongo-tls-insecure", false, "Enable TLS insecure")
	MongoCAFile      = FlagSet.String("mongo-tls-ca-file", "", "Mongo TLS CA file")
	MongoCertFile    = FlagSet.String("mongo-tls-cert-file", "", "Mongo TLS cert file")
	MongoKeyFile     = FlagSet.String("mongo-tls-key-file", "", "Mongo TLS key file")
)

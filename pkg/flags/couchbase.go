package flags

var (
	CouchbaseAddress    = FlagSet.String("couchbase-address", "", "Couchbase address")
	CouchbaseUser       = FlagSet.String("couchbase-user", "", "Couchbase user")
	CouchbasePassword   = FlagSet.String("couchbase-password", "", "Couchbase password")
	CouchbaseBucketName = FlagSet.String("couchbase-bucket", "", "Couchbase bucket name")
	CouchbaseScope      = FlagSet.String("couchbase-scope", "_default", "Couchbase scope")
	CouchbaseCollection = FlagSet.String("couchbase-collection", "_default", "Couchbase collection")
	CouchbaseID         = FlagSet.String("couchbase-id", "", "Couchbase id")
	// TLS
	CouchbaseEnableTLS   = FlagSet.Bool("couchbase-enable-tls", false, "Enable TLS")
	CouchbaseTLSInsecure = FlagSet.Bool("couchbase-tls-insecure", false, "Enable TLS insecure")
	CouchbaseCAFile      = FlagSet.String("couchbase-tls-ca-file", "", "Couchbase TLS CA file")
	CouchbaseCertFile    = FlagSet.String("couchbase-tls-cert-file", "", "Couchbase TLS cert file")
	CouchbaseKeyFile     = FlagSet.String("couchbase-tls-key-file", "", "Couchbase TLS key file")
)

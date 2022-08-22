package flags

var (
	CockroachDBHost        = FlagSet.String("cockroach-host", "", "CockroachDB host")
	CockroachDBPort        = FlagSet.String("cockroach-port", "26257", "CockroachDB port")
	CockroachDBUser        = FlagSet.String("cockroach-user", "", "CockroachDB user")
	CockroachDBPassword    = FlagSet.String("cockroach-password", "", "CockroachDB password")
	CockroachDBDatabase    = FlagSet.String("cockroach-database", "", "CockroachDB database")
	CockroachDBSSLMode     = FlagSet.String("cockroach-ssl-mode", "disable", "CockroachDB SSL mode")
	CockroachDBQuery       = FlagSet.String("cockroach-query", "", "CockroachDB query")
	CockroachDBQueryParams = FlagSet.String("cockroach-params", "", "CockroachDB query params")
	CockroachDBRoutingID   = FlagSet.String("cockroach-routing-id", "", "CockroachDB routing id")
	CockroachDBTLSRootCert = FlagSet.String("cockroach-tls-root-cert", "", "CockroachDB TLS root cert")
	CockroachDBTLSCert     = FlagSet.String("cockroach-tls-cert", "", "CockroachDB TLS cert")
	CockroachDBTLSKey      = FlagSet.String("cockroach-tls-key", "", "CockroachDB TLS key")
)

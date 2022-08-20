package flags

var (
	PsqlHost        = FlagSet.String("psql-host", "", "PostgreSQL host")
	PsqlPort        = FlagSet.String("psql-port", "5432", "PostgreSQL port")
	PsqlUser        = FlagSet.String("psql-user", "", "PostgreSQL user")
	PsqlPassword    = FlagSet.String("psql-password", "", "PostgreSQL password")
	PsqlDatabase    = FlagSet.String("psql-database", "", "PostgreSQL database")
	PsqlSSLMode     = FlagSet.String("psql-ssl-mode", "disable", "PostgreSQL SSL mode")
	PsqlQuery       = FlagSet.String("psql-query", "", "PostgreSQL query")
	PsqlQueryParams = FlagSet.String("psql-params", "", "PostgreSQL query params")
)

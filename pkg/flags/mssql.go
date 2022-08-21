package flags

var (
	MSSqlHost        = FlagSet.String("mssql-host", "", "MySQL host")
	MSSqlPort        = FlagSet.String("mssql-port", "1433", "MySQL port")
	MSSqlUser        = FlagSet.String("mssql-user", "", "MySQL user")
	MSSqlPassword    = FlagSet.String("mssql-password", "", "MySQL password")
	MSSqlDatabase    = FlagSet.String("mssql-database", "", "MySQL database")
	MSSqlQuery       = FlagSet.String("mssql-query", "", "MySQL query")
	MSSqlQueryParams = FlagSet.String("mssql-params", "", "MySQL query params")
)

package flags

var (
	MysqlHost        = FlagSet.String("mysql-host", "", "MySQL host")
	MysqlPort        = FlagSet.String("mysql-port", "3306", "MySQL port")
	MysqlUser        = FlagSet.String("mysql-user", "", "MySQL user")
	MysqlPassword    = FlagSet.String("mysql-password", "", "MySQL password")
	MysqlDatabase    = FlagSet.String("mysql-database", "", "MySQL database")
	MysqlQuery       = FlagSet.String("mysql-query", "", "MySQL query")
	MysqlQueryParams = FlagSet.String("mysql-params", "", "MySQL query params")
)

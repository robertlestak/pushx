package flags

var (
	CassandraHosts       = FlagSet.String("cassandra-hosts", "", "Cassandra hosts")
	CassandraUser        = FlagSet.String("cassandra-user", "", "Cassandra user")
	CassandraPassword    = FlagSet.String("cassandra-password", "", "Cassandra password")
	CassandraKeyspace    = FlagSet.String("cassandra-keyspace", "", "Cassandra keyspace")
	CassandraConsistency = FlagSet.String("cassandra-consistency", "QUORUM", "Cassandra consistency")
	CassandraQuery       = FlagSet.String("cassandra-query", "", "Cassandra query")
	CassandraQueryParams = FlagSet.String("cassandra-params", "", "Cassandra query params")
)

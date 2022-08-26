package flags

var (
	ScyllaHosts       = FlagSet.String("scylla-hosts", "", "Scylla hosts")
	ScyllaUser        = FlagSet.String("scylla-user", "", "Scylla user")
	ScyllaPassword    = FlagSet.String("scylla-password", "", "Scylla password")
	ScyllaKeyspace    = FlagSet.String("scylla-keyspace", "", "Scylla keyspace")
	ScyllaConsistency = FlagSet.String("scylla-consistency", "QUORUM", "Scylla consistency")
	ScyllaLocalDC     = FlagSet.String("scylla-local-dc", "", "Scylla local dc")
	ScyllaQuery       = FlagSet.String("scylla-query", "", "Scylla query")
	ScyllaQueryParams = FlagSet.String("scylla-params", "", "Scylla query params")
)

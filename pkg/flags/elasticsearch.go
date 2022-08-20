package flags

var (
	ElasticsearchAddress       = FlagSet.String("elasticsearch-address", "", "Elasticsearch address")
	ElasticsearchUsername      = FlagSet.String("elasticsearch-username", "", "Elasticsearch username")
	ElasticsearchPassword      = FlagSet.String("elasticsearch-password", "", "Elasticsearch password")
	ElasticsearchTLSSkipVerify = FlagSet.Bool("elasticsearch-tls-skip-verify", false, "Elasticsearch TLS skip verify")
	ElasticsearchEnableTLS     = FlagSet.Bool("elasticsearch-enable-tls", false, "Elasticsearch enable TLS")
	ElasticsearchCAFile        = FlagSet.String("elasticsearch-tls-ca-file", "", "Elasticsearch TLS CA file")
	ElasticsearchCertFile      = FlagSet.String("elasticsearch-tls-cert-file", "", "Elasticsearch TLS cert file")
	ElasticsearchKeyFile       = FlagSet.String("elasticsearch-tls-key-file", "", "Elasticsearch TLS key file")
	ElasticsearchIndex         = FlagSet.String("elasticsearch-index", "", "Elasticsearch index")
	ElasticsearchDocID         = FlagSet.String("elasticsearch-doc-id", "", "Elasticsearch doc id")
)

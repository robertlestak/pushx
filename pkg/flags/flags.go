package flags

import "flag"

var (
	FlagSet   = flag.NewFlagSet("pushx", flag.ContinueOnError)
	Driver    = FlagSet.String("driver", "", "driver to use. (activemq, aws-dynamo, aws-s3, aws-sqs, cassandra, centauri, cockroach, couchbase, elasticsearch, fs, gcp-bq, gcp-firestore, gcp-gcs, gcp-pubsub, github, http, kafka, local, mongodb, mssql, mysql, nats, nfs, nsq, postgres, pulsar, rabbitmq, redis-list, redis-pubsub, redis-stream)")
	InputFile = FlagSet.String("in-file", "-", "input file to use. (default: stdin)")
	InputStr  = FlagSet.String("in", "", "input string to use. Will take precedence over -in-file")
	Output    = FlagSet.String("out", "", "output file to use in addition to the driver. If '-' then stdout is used.")
)

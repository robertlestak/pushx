package drivers

import (
	"errors"

	"github.com/robertlestak/pushx/drivers/activemq"
	"github.com/robertlestak/pushx/drivers/aws"
	"github.com/robertlestak/pushx/drivers/cassandra"
	"github.com/robertlestak/pushx/drivers/centauri"
	"github.com/robertlestak/pushx/drivers/elasticsearch"
	"github.com/robertlestak/pushx/drivers/fs"
	"github.com/robertlestak/pushx/drivers/gcp"
	"github.com/robertlestak/pushx/drivers/http"
	"github.com/robertlestak/pushx/drivers/kafka"
	"github.com/robertlestak/pushx/drivers/local"
	"github.com/robertlestak/pushx/drivers/mongodb"
	"github.com/robertlestak/pushx/drivers/mysql"
	"github.com/robertlestak/pushx/drivers/nats"
	"github.com/robertlestak/pushx/drivers/nfs"
	"github.com/robertlestak/pushx/drivers/nsq"
	"github.com/robertlestak/pushx/drivers/postgres"
	"github.com/robertlestak/pushx/drivers/pulsar"
	"github.com/robertlestak/pushx/drivers/rabbitmq"
	"github.com/robertlestak/pushx/drivers/redis"
)

type DriverName string

var (
	ActiveMQ          DriverName = "activemq"
	AWSS3             DriverName = "aws-s3"
	AWSSQS            DriverName = "aws-sqs"
	AWSDynamoDB       DriverName = "aws-dynamo"
	CassandraDB       DriverName = "cassandra"
	Centauri          DriverName = "centauri"
	Elasticsearch     DriverName = "elasticsearch"
	FS                DriverName = "fs"
	HTTP              DriverName = "http"
	Kafka             DriverName = "kafka"
	GCPBQ             DriverName = "gcp-bq"
	GCPFirestore      DriverName = "gcp-firestore"
	GCPGCS            DriverName = "gcp-gcs"
	GCPPubSub         DriverName = "gcp-pubsub"
	MongoDB           DriverName = "mongodb"
	MySQL             DriverName = "mysql"
	Nats              DriverName = "nats"
	NSQ               DriverName = "nsq"
	NFS               DriverName = "nfs"
	Postgres          DriverName = "postgres"
	Pulsar            DriverName = "pulsar"
	Rabbit            DriverName = "rabbitmq"
	RedisList         DriverName = "redis-list"
	RedisSubscription DriverName = "redis-pubsub"
	RedisStream       DriverName = "redis-stream"
	Local             DriverName = "local"
	ErrDriverNotFound            = errors.New("driver not found")
)

// Get returns the driver with the given name.
func GetDriver(name DriverName) Driver {
	switch name {
	case ActiveMQ:
		return &activemq.ActiveMQ{}
	case AWSS3:
		return &aws.S3{}
	case AWSSQS:
		return &aws.SQS{}
	case AWSDynamoDB:
		return &aws.Dynamo{}
	case CassandraDB:
		return &cassandra.Cassandra{}
	case Centauri:
		return &centauri.Centauri{}
	case Elasticsearch:
		return &elasticsearch.Elasticsearch{}
	case FS:
		return &fs.FS{}
	case GCPBQ:
		return &gcp.BQ{}
	case GCPFirestore:
		return &gcp.GCPFirestore{}
	case GCPGCS:
		return &gcp.GCS{}
	case GCPPubSub:
		return &gcp.GCPPubSub{}
	case HTTP:
		return &http.HTTP{}
	case Kafka:
		return &kafka.Kafka{}
	case MongoDB:
		return &mongodb.Mongo{}
	case MySQL:
		return &mysql.Mysql{}
	case Nats:
		return &nats.NATS{}
	case NSQ:
		return &nsq.NSQ{}
	case NFS:
		return &nfs.NFS{}
	case Postgres:
		return &postgres.Postgres{}
	case Pulsar:
		return &pulsar.Pulsar{}
	case Rabbit:
		return &rabbitmq.RabbitMQ{}
	case RedisList:
		return &redis.RedisList{}
	case RedisSubscription:
		return &redis.RedisPubSub{}
	case RedisStream:
		return &redis.RedisStream{}
	case Local:
		return &local.Local{}
	}
	return nil
}

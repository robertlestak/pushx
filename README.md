# pushx - cloud agnostic data push

pushx is a cloud agnostic abstraction binary for pushing data to various persistence layers.

pushx is a single compiled binary that can be packaged in your existing job code container, configured with environment variables or command line flags, and included in any existing data pipelines and workflows.

Used in conjunction with [procx](https://github.com/robertlestak/procx), pushx can be used to build entirely cloud and provider agnostic data pipelines.

## Execution

pushx is configured with either environment variables or a set of command line flags, and accepts data input either from stdin, a file, or command line arguments.

```bash
# cli args
echo -n hello world | pushx -driver redis-list ...
# or, env vars
export PUSHX_DRIVER=redis-list
...
echo -n hello world | pushx
```

### Payload

By default, pushx will read input data from stdin. If `-in-file` is provided, pushx will read input data from the specified file, and if `-in` is provided, pushx will read input data from the specified command line argument.

### Pipelining

By default, pushx will consume the input data, send it to the configured data provider, and exit with a 0 status code on success and a non-zero exit code on failure to send to the data provider. However if you would like to include `pushx` in a larger shell pipeline, you can use `-out=-` (or `-out=filename.txt`) to pipe the input to the output to pass on to the next command. This also allows you to send the same data to multiple data providers.

```bash
echo -n hello world | pushx -driver redis-list -out=- | pushx -driver gcp-pubsub -out=- | pushx -driver gcp-bq
```

#### Relational Driver JSON Parsing

For drivers which are non-structured (ex. `fs`, `aws-s3`, `redis-list`, etc.), pushx will send the input data as-is to the driver. However for drivers which enforce some relational schema such as SQL-based drivers, you will need to provide an input query which will be executed to insert the input data. You can provide a `{{pushx_payload}}` placeholder in your query / parameters which will be replaced with the entire input data. For example:

```bash
echo 'the data' | pushx -driver postgres \
    ...
    -psql-query "INSERT INTO table (data) VALUES ($1)"
    -psql-params "{{pushx_payload}}"
```

However if your input data is a JSON object, you may want to convert this to a relational format when inserting into your column-oriented database. You can use `{{mustache}}` syntax to extract specific fields from the input data and insert them into your query. For example:

```bash
echo '{"id": 1, "name": "John"}' | pushx -driver postgres \
    ...
    -psql-query "INSERT INTO table (id, name) VALUES ($1, $2)"
    -psql-params "{{id}},{{name}}"
```

This also supports deeply nested fields, using `gjson` syntax.

```bash
echo '{"id": 1, "name": "John", "address": {"street": "123 Main St", "city": "Anytown"}}' | pushx -driver postgres \
    ...
    -psql-query "INSERT INTO table (id, name, street, city) VALUES ($1, $2, $3, $4)"
    -psql-params "{{id}},{{name}},{{address.street}},{{address.city}}"
```

## Drivers

Currently, the following drivers are supported:

- [ActiveMQ](#activemq) (`activemq`)
- [AWS DynamoDB](#aws-dynamodb) (`aws-dynamo`)
- [AWS S3](#aws-s3) (`aws-s3`)
- [AWS SQS](#aws-sqs) (`aws-sqs`)
- [Cassandra](#cassandra) (`cassandra`)
- [Centauri](#centauri) (`centauri`)
- [Cockroach](#cockroach) (`cockroach`)
- [Couchbase](#couchbase) (`couchbase`)
- [Elasticsearch](#elasticsearch) (`elasticsearch`)
- [FS](#fs) (`fs`)
- [GCP Big Query](#gcp-bq) (`gcp-bq`)
- [GCP Cloud Storage](#gcp-gcs) (`gcp-gcs`)
- [GCP Firestore](#gcp-firestore) (`gcp-firestore`)
- [GCP Pub/Sub](#gcp-pubsub) (`gcp-pubsub`)
- [GitHub](#github) (`github`)
- [HTTP](#http) (`http`)
- [Kafka](#kafka) (`kafka`)
- [PostgreSQL](#postgresql) (`postgres`)
- [Pulsar](#pulsar) (`pulsar`)
- [MongoDB](#mongodb) (`mongodb`)
- [MSSQL](#mssql) (`mssql`)
- [MySQL](#mysql) (`mysql`)
- [NATS](#nats) (`nats`)
- [NFS](#nfs) (`nfs`)
- [NSQ](#nsq) (`nsq`)
- [RabbitMQ](#rabbitmq) (`rabbitmq`)
- [Redis List](#redis-list) (`redis-list`)
- [Redis Pub/Sub](#redis-pubsub) (`redis-pubsub`)
- [Redis Stream](#redis-stream) (`redis-stream`)
- [Local](#local) (`local`)

Plans to add more drivers in the future, and PRs are welcome.

See [Driver Examples](#driver-examples) for more information.

## Install

```bash
curl -SsL https://raw.githubusercontent.com/robertlestak/pushx/main/scripts/install.sh | bash -e
```

### A note on permissions

Depending on the path of `INSTALL_DIR` and the permissions of the user running the installation script, you may get a Permission Denied error if you are trying to move the binary into a location which your current user does not have access to. This is most often the case when running the script as a non-root user yet trying to install into `/usr/local/bin`. To fix this, you can either:

Create a `$HOME/bin` directory in your current user home directory. This will be the default installation directory. Be sure to add this to your `$PATH` environment variable.

Use `sudo` to run the installation script, to install into `/usr/local/bin` 

```bash
curl -SsL https://raw.githubusercontent.com/robertlestak/pushx/main/scripts/install.sh | sudo bash -e
```

### Build From Source

```bash
mkdir -p bin
go build -o bin/pushx cmd/pushx/*.go
```

#### Building for a Specific Driver

By default, the `pushx` binary is compiled for all drivers. This is to enable a truly build-once-run-anywhere experience. However some users may want a smaller binary for embedded workloads. To enable this, you can run `make listdrivers` to get the full list of available drivers, and `make slim drivers="driver1 driver2 driver3 ..."` - listing each driver separated by a space - to build a slim binary with just the specified driver(s).

While building for a specific driver may seem contrary to the ethos of pushx, the decoupling between the job queue and work still enables a write-once-run-anywhere experience, and simply requires DevOps to rebuild the image with your new drivers if you are shifting upstream data sources.

## Usage

```bash
Usage: pushx [options]
  -activemq-address string
    	ActiveMQ STOMP address
  -activemq-enable-tls
    	Enable TLS
  -activemq-name string
    	ActiveMQ name
  -activemq-tls-ca-file string
    	TLS CA
  -activemq-tls-cert-file string
    	TLS cert
  -activemq-tls-insecure
    	Enable TLS insecure
  -activemq-tls-key-file string
    	TLS key
  -aws-dynamo-table string
    	AWS DynamoDB table name
  -aws-load-config
    	load AWS config from ~/.aws/config
  -aws-region string
    	AWS region
  -aws-role-arn string
    	AWS role ARN
  -aws-s3-acl string
    	AWS S3 ACL
  -aws-s3-bucket string
    	AWS S3 bucket
  -aws-s3-key string
    	AWS S3 key
  -aws-s3-tags string
    	AWS S3 tags. Comma separated list of key=value pairs
  -aws-sqs-queue-url string
    	AWS SQS queue URL
  -cassandra-consistency string
    	Cassandra consistency (default "QUORUM")
  -cassandra-hosts string
    	Cassandra hosts
  -cassandra-keyspace string
    	Cassandra keyspace
  -cassandra-params string
    	Cassandra query params
  -cassandra-password string
    	Cassandra password
  -cassandra-query string
    	Cassandra query
  -cassandra-user string
    	Cassandra user
  -centauri-channel string
    	Centauri channel (default "default")
  -centauri-filename string
    	Centauri filename
  -centauri-message-type string
    	Centauri message type. One of: bytes, file (default "bytes")
  -centauri-peer-url string
    	Centauri peer URL
  -centauri-public-key string
    	Centauri public key
  -centauri-public-key-base64 string
    	Centauri public key base64
  -cockroach-database string
    	CockroachDB database
  -cockroach-host string
    	CockroachDB host
  -cockroach-params string
    	CockroachDB query params
  -cockroach-password string
    	CockroachDB password
  -cockroach-port string
    	CockroachDB port (default "26257")
  -cockroach-query string
    	CockroachDB query
  -cockroach-routing-id string
    	CockroachDB routing id
  -cockroach-ssl-mode string
    	CockroachDB SSL mode (default "disable")
  -cockroach-tls-cert string
    	CockroachDB TLS cert
  -cockroach-tls-key string
    	CockroachDB TLS key
  -cockroach-tls-root-cert string
    	CockroachDB TLS root cert
  -cockroach-user string
    	CockroachDB user
  -couchbase-address string
    	Couchbase address
  -couchbase-bucket string
    	Couchbase bucket name
  -couchbase-collection string
    	Couchbase collection (default "_default")
  -couchbase-enable-tls
    	Enable TLS
  -couchbase-id string
    	Couchbase id
  -couchbase-password string
    	Couchbase password
  -couchbase-scope string
    	Couchbase scope (default "_default")
  -couchbase-tls-ca-file string
    	Couchbase TLS CA file
  -couchbase-tls-cert-file string
    	Couchbase TLS cert file
  -couchbase-tls-insecure
    	Enable TLS insecure
  -couchbase-tls-key-file string
    	Couchbase TLS key file
  -couchbase-user string
    	Couchbase user
  -driver string
    	driver to use. (activemq, aws-dynamo, aws-s3, aws-sqs, cassandra, centauri, cockroach, couchbase, elasticsearch, fs, gcp-bq, gcp-firestore, gcp-gcs, gcp-pubsub, github, http, kafka, local, mongodb, mssql, mysql, nats, nfs, nsq, postgres, pulsar, rabbitmq, redis-list, redis-pubsub, redis-stream)
  -elasticsearch-address string
    	Elasticsearch address
  -elasticsearch-doc-id string
    	Elasticsearch doc id
  -elasticsearch-enable-tls
    	Elasticsearch enable TLS
  -elasticsearch-index string
    	Elasticsearch index
  -elasticsearch-password string
    	Elasticsearch password
  -elasticsearch-tls-ca-file string
    	Elasticsearch TLS CA file
  -elasticsearch-tls-cert-file string
    	Elasticsearch TLS cert file
  -elasticsearch-tls-key-file string
    	Elasticsearch TLS key file
  -elasticsearch-tls-skip-verify
    	Elasticsearch TLS skip verify
  -elasticsearch-username string
    	Elasticsearch username
  -fs-folder string
    	FS folder
  -fs-key string
    	FS key
  -gcp-bq-query string
    	GCP BigQuery query
  -gcp-firestore-collection string
    	GCP Firestore collection
  -gcp-firestore-id string
    	GCP Firestore document ID. If empty, a new document ID will be created
  -gcp-gcs-bucket string
    	GCP GCS bucket
  -gcp-gcs-key string
    	GCP GCS key
  -gcp-project-id string
    	GCP project ID
  -gcp-pubsub-topic string
    	GCP Pub/Sub topic name
  -github-base-branch string
    	base branch for PR
  -github-branch string
    	branch for PR.
  -github-commit-email string
    	commit email
  -github-commit-message string
    	commit message
  -github-commit-name string
    	commit name
  -github-file string
    	GitHub file
  -github-open-pr
    	open PR on changes. Default: false
  -github-owner string
    	GitHub owner
  -github-pr-body string
    	PR body
  -github-pr-title string
    	PR title
  -github-ref string
    	GitHub ref
  -github-repo string
    	GitHub repo
  -github-token string
    	GitHub token
  -http-content-type string
    	HTTP content type
  -http-enable-tls
    	HTTP enable tls
  -http-headers string
    	HTTP headers
  -http-method string
    	HTTP method (default "POST")
  -http-successful-status-codes string
    	HTTP successful status codes
  -http-tls-ca-file string
    	HTTP tls ca file
  -http-tls-cert-file string
    	HTTP tls cert file
  -http-tls-insecure
    	HTTP tls insecure
  -http-tls-key-file string
    	HTTP tls key file
  -http-url string
    	HTTP url
  -in string
    	input string to use. Will take precedence over -in-file
  -in-file string
    	input file to use. (default: stdin) (default "-")
  -kafka-brokers string
    	Kafka brokers, comma separated
  -kafka-enable-sasl
    	Enable SASL
  -kafka-enable-tls
    	Enable TLS
  -kafka-sasl-password string
    	Kafka SASL password
  -kafka-sasl-type string
    	Kafka SASL type. Can be either 'scram' or 'plain'
  -kafka-sasl-username string
    	Kafka SASL user
  -kafka-tls-ca-file string
    	Kafka TLS CA file
  -kafka-tls-cert-file string
    	Kafka TLS cert file
  -kafka-tls-insecure
    	Enable TLS insecure
  -kafka-tls-key-file string
    	Kafka TLS key file
  -kafka-topic string
    	Kafka topic
  -mongo-auth-source string
    	MongoDB auth source
  -mongo-collection string
    	MongoDB collection
  -mongo-database string
    	MongoDB database
  -mongo-enable-tls
    	Enable TLS
  -mongo-host string
    	MongoDB host
  -mongo-password string
    	MongoDB password
  -mongo-port string
    	MongoDB port (default "27017")
  -mongo-tls-ca-file string
    	Mongo TLS CA file
  -mongo-tls-cert-file string
    	Mongo TLS cert file
  -mongo-tls-insecure
    	Enable TLS insecure
  -mongo-tls-key-file string
    	Mongo TLS key file
  -mongo-user string
    	MongoDB user
  -mssql-database string
    	MySQL database
  -mssql-host string
    	MySQL host
  -mssql-params string
    	MySQL query params
  -mssql-password string
    	MySQL password
  -mssql-port string
    	MySQL port (default "1433")
  -mssql-query string
    	MySQL query
  -mssql-user string
    	MySQL user
  -mysql-database string
    	MySQL database
  -mysql-host string
    	MySQL host
  -mysql-params string
    	MySQL query params
  -mysql-password string
    	MySQL password
  -mysql-port string
    	MySQL port (default "3306")
  -mysql-query string
    	MySQL query
  -mysql-user string
    	MySQL user
  -nats-creds-file string
    	NATS creds file
  -nats-enable-tls
    	NATS enable TLS
  -nats-jwt-file string
    	NATS JWT file
  -nats-nkey-file string
    	NATS NKey file
  -nats-password string
    	NATS password
  -nats-subject string
    	NATS subject
  -nats-tls-ca-file string
    	NATS TLS CA file
  -nats-tls-cert-file string
    	NATS TLS cert file
  -nats-tls-insecure
    	NATS TLS insecure
  -nats-tls-key-file string
    	NATS TLS key file
  -nats-token string
    	NATS token
  -nats-url string
    	NATS URL
  -nats-username string
    	NATS username
  -nfs-folder string
    	NFS folder
  -nfs-host string
    	NFS host
  -nfs-key string
    	NFS key
  -nfs-target string
    	NFS target
  -nsq-enable-tls
    	Enable TLS
  -nsq-nsqd-address string
    	NSQ nsqd address
  -nsq-nsqlookupd-address string
    	NSQ nsqlookupd address
  -nsq-tls-ca-file string
    	NSQ TLS CA file
  -nsq-tls-cert-file string
    	NSQ TLS cert file
  -nsq-tls-key-file string
    	NSQ TLS key file
  -nsq-tls-skip-verify
    	NSQ TLS skip verify
  -nsq-topic string
    	NSQ topic
  -out string
    	output file to use in addition to the driver. If '-' then stdout is used.
  -psql-database string
    	PostgreSQL database
  -psql-host string
    	PostgreSQL host
  -psql-params string
    	PostgreSQL query params
  -psql-password string
    	PostgreSQL password
  -psql-port string
    	PostgreSQL port (default "5432")
  -psql-query string
    	PostgreSQL query
  -psql-ssl-mode string
    	PostgreSQL SSL mode (default "disable")
  -psql-tls-cert string
    	PostgreSQL TLS cert
  -psql-tls-key string
    	PostgreSQL TLS key
  -psql-tls-root-cert string
    	PostgreSQL TLS root cert
  -psql-user string
    	PostgreSQL user
  -pulsar-address string
    	Pulsar address
  -pulsar-auth-cert-file string
    	Pulsar auth cert file
  -pulsar-auth-key-file string
    	Pulsar auth key file
  -pulsar-auth-oauth-params string
    	Pulsar auth oauth params
  -pulsar-auth-token string
    	Pulsar auth token
  -pulsar-auth-token-file string
    	Pulsar auth token file
  -pulsar-producer-name string
    	Pulsar producer name
  -pulsar-tls-allow-insecure-connection
    	Pulsar TLS allow insecure connection
  -pulsar-tls-trust-certs-file string
    	Pulsar TLS trust certs file path
  -pulsar-tls-validate-hostname
    	Pulsar TLS validate hostname
  -pulsar-topic string
    	Pulsar topic
  -rabbitmq-exchange string
    	RabbitMQ exchange
  -rabbitmq-queue string
    	RabbitMQ queue
  -rabbitmq-url string
    	RabbitMQ URL
  -redis-enable-tls
    	Enable TLS
  -redis-host string
    	Redis host
  -redis-key string
    	Redis key
  -redis-message-id string
    	Redis stream message id (default "*")
  -redis-password string
    	Redis password
  -redis-port string
    	Redis port (default "6379")
  -redis-tls-ca-file string
    	Redis TLS CA file
  -redis-tls-cert-file string
    	Redis TLS cert file
  -redis-tls-key-file string
    	Redis TLS key file
  -redis-tls-skip-verify
    	Redis TLS skip verify
```

### Environment Variables

- `AWS_REGION`
- `AWS_SDK_LOAD_CONFIG`
- `LOG_LEVEL`
- `NSQ_LOG_LEVEL`
- `PUSHX_ACTIVEMQ_ADDRESS`
- `PUSHX_ACTIVEMQ_ENABLE_TLS`
- `PUSHX_ACTIVEMQ_NAME`
- `PUSHX_ACTIVEMQ_TLS_CA_FILE`
- `PUSHX_ACTIVEMQ_TLS_CERT_FILE`
- `PUSHX_ACTIVEMQ_TLS_INSECURE`
- `PUSHX_ACTIVEMQ_TLS_KEY_FILE`
- `PUSHX_AWS_DYNAMO_TABLE`
- `PUSHX_AWS_LOAD_CONFIG`
- `PUSHX_AWS_REGION`
- `PUSHX_AWS_ROLE_ARN`
- `PUSHX_AWS_S3_ACL`
- `PUSHX_AWS_S3_BUCKET`
- `PUSHX_AWS_S3_KEY`
- `PUSHX_AWS_S3_TAGS`
- `PUSHX_AWS_SQS_QUEUE_URL`
- `PUSHX_AWS_SQS_ROLE_ARN`
- `PUSHX_CASSANDRA_CONSISTENCY`
- `PUSHX_CASSANDRA_HOSTS`
- `PUSHX_CASSANDRA_KEYSPACE`
- `PUSHX_CASSANDRA_PARAMS`
- `PUSHX_CASSANDRA_PASSWORD`
- `PUSHX_CASSANDRA_QUERY`
- `PUSHX_CASSANDRA_USER`
- `PUSHX_CENTAURI_CHANNEL`
- `PUSHX_CENTAURI_FILENAME`
- `PUSHX_CENTAURI_MESSAGE_TYPE`
- `PUSHX_CENTAURI_PEER_URL`
- `PUSHX_CENTAURI_PUBLIC_KEY`
- `PUSHX_CENTAURI_PUBLIC_KEY_BASE64`
- `PUSHX_COCKROACH_DATABASE`
- `PUSHX_COCKROACH_HOST`
- `PUSHX_COCKROACH_PASSWORD`
- `PUSHX_COCKROACH_PORT`
- `PUSHX_COCKROACH_QUERY`
- `PUSHX_COCKROACH_QUERY_PARAMS`
- `PUSHX_COCKROACH_ROUTING_ID`
- `PUSHX_COCKROACH_SSL_MODE`
- `PUSHX_COCKROACH_TLS_CERT`
- `PUSHX_COCKROACH_TLS_KEY`
- `PUSHX_COCKROACH_TLS_ROOT_CERT`
- `PUSHX_COCKROACH_USER`
- `PUSHX_COUCHBASE_BUCKET_NAME`
- `PUSHX_COUCHBASE_COLLECTION`
- `PUSHX_COUCHBASE_ENABLE_TLS`
- `PUSHX_COUCHBASE_ID`
- `PUSHX_COUCHBASE_PASSWORD`
- `PUSHX_COUCHBASE_SCOPE`
- `PUSHX_COUCHBASE_TLS_CA_FILE`
- `PUSHX_COUCHBASE_TLS_CERT_FILE`
- `PUSHX_COUCHBASE_TLS_INSECURE`
- `PUSHX_COUCHBASE_TLS_KEY_FILE`
- `PUSHX_COUCHBASE_USER`
- `PUSHX_DRIVER`
- `PUSHX_ELASTICSEARCH_ADDRESS`
- `PUSHX_ELASTICSEARCH_DOC_ID`
- `PUSHX_ELASTICSEARCH_ENABLE_TLS`
- `PUSHX_ELASTICSEARCH_INDEX`
- `PUSHX_ELASTICSEARCH_PASSWORD`
- `PUSHX_ELASTICSEARCH_TLS_CA_FILE`
- `PUSHX_ELASTICSEARCH_TLS_CERT_FILE`
- `PUSHX_ELASTICSEARCH_TLS_KEY_FILE`
- `PUSHX_ELASTICSEARCH_TLS_SKIP_VERIFY`
- `PUSHX_ELASTICSEARCH_USERNAME`
- `PUSHX_FS_FOLDER`
- `PUSHX_FS_KEY`
- `PUSHX_GCP_BQ_QUERY`
- `PUSHX_GCP_FIRESTORE_COLLECTION`
- `PUSHX_GCP_FIRESTORE_ID`
- `PUSHX_GCP_GCS_BUCKET`
- `PUSHX_GCP_GCS_KEY`
- `PUSHX_GCP_PROJECT_ID`
- `PUSHX_GCP_TOPIC`
- `PUSHX_GITHUB_BASE_BRANCH`
- `PUSHX_GITHUB_BRANCH`
- `PUSHX_GITHUB_COMMIT_EMAIL`
- `PUSHX_GITHUB_COMMIT_MESSAGE`
- `PUSHX_GITHUB_COMMIT_NAME`
- `PUSHX_GITHUB_FILE`
- `PUSHX_GITHUB_OPEN_PR`
- `PUSHX_GITHUB_OWNER`
- `PUSHX_GITHUB_PR_BODY`
- `PUSHX_GITHUB_PR_TITLE`
- `PUSHX_GITHUB_REF`
- `PUSHX_GITHUB_REPO`
- `PUSHX_GITHUB_TOKEN`
- `PUSHX_HTTP_ENABLE_TLS`
- `PUSHX_HTTP_REQUEST_CONTENT_TYPE`
- `PUSHX_HTTP_REQUEST_HEADERS`
- `PUSHX_HTTP_REQUEST_METHOD`
- `PUSHX_HTTP_REQUEST_SUCCESSFUL_STATUS_CODES`
- `PUSHX_HTTP_REQUEST_URL`
- `PUSHX_HTTP_TLS_CA_FILE`
- `PUSHX_HTTP_TLS_CERT_FILE`
- `PUSHX_HTTP_TLS_KEY_FILE`
- `PUSHX_INPUT_FILE`
- `PUSHX_INPUT_STR`
- `PUSHX_KAFKA_BROKERS`
- `PUSHX_KAFKA_ENABLE_SASL`
- `PUSHX_KAFKA_ENABLE_TLS`
- `PUSHX_KAFKA_SASL_PASSWORD`
- `PUSHX_KAFKA_SASL_TYPE`
- `PUSHX_KAFKA_SASL_USERNAME`
- `PUSHX_KAFKA_TLS_CA_FILE`
- `PUSHX_KAFKA_TLS_CERT_FILE`
- `PUSHX_KAFKA_TLS_INSECURE`
- `PUSHX_KAFKA_TLS_KEY_FILE`
- `PUSHX_KAFKA_TOPIC`
- `PUSHX_MONGO_AUTH_SOURCE`
- `PUSHX_MONGO_COLLECTION`
- `PUSHX_MONGO_DATABASE`
- `PUSHX_MONGO_ENABLE_TLS`
- `PUSHX_MONGO_HOST`
- `PUSHX_MONGO_PASSWORD`
- `PUSHX_MONGO_PORT`
- `PUSHX_MONGO_TLS_CA_FILE`
- `PUSHX_MONGO_TLS_CERT_FILE`
- `PUSHX_MONGO_TLS_INSECURE`
- `PUSHX_MONGO_TLS_KEY_FILE`
- `PUSHX_MONGO_USER`
- `PUSHX_MSSQL_DATABASE`
- `PUSHX_MSSQL_HOST`
- `PUSHX_MSSQL_PASSWORD`
- `PUSHX_MSSQL_PORT`
- `PUSHX_MSSQL_QUERY`
- `PUSHX_MSSQL_QUERY_PARAMS`
- `PUSHX_MSSQL_USER`
- `PUSHX_MYSQL_DATABASE`
- `PUSHX_MYSQL_HOST`
- `PUSHX_MYSQL_PASSWORD`
- `PUSHX_MYSQL_PORT`
- `PUSHX_MYSQL_QUERY`
- `PUSHX_MYSQL_QUERY_PARAMS`
- `PUSHX_MYSQL_USER`
- `PUSHX_NATS_CREDS_FILE`
- `PUSHX_NATS_ENABLE_TLS`
- `PUSHX_NATS_JWT_FILE`
- `PUSHX_NATS_NKEY_FILE`
- `PUSHX_NATS_PASSWORD`
- `PUSHX_NATS_SUBJECT`
- `PUSHX_NATS_TLS_CA_FILE`
- `PUSHX_NATS_TLS_CERT_FILE`
- `PUSHX_NATS_TLS_INSECURE`
- `PUSHX_NATS_TLS_KEY_FILE`
- `PUSHX_NATS_TOKEN`
- `PUSHX_NATS_URL`
- `PUSHX_NATS_USERNAME`
- `PUSHX_NFS_FOLDER`
- `PUSHX_NFS_HOST`
- `PUSHX_NFS_KEY`
- `PUSHX_NFS_TARGET`
- `PUSHX_NSQ_ENABLE_TLS`
- `PUSHX_NSQ_NSQD_ADDRESS`
- `PUSHX_NSQ_NSQLOOKUPD_ADDRESS`
- `PUSHX_NSQ_TLS_CA_FILE`
- `PUSHX_NSQ_TLS_CERT_FILE`
- `PUSHX_NSQ_TLS_INSECURE`
- `PUSHX_NSQ_TLS_KEY_FILE`
- `PUSHX_NSQ_TOPIC`
- `PUSHX_OUTPUT`
- `PUSHX_PSQL_DATABASE`
- `PUSHX_PSQL_HOST`
- `PUSHX_PSQL_PASSWORD`
- `PUSHX_PSQL_PORT`
- `PUSHX_PSQL_QUERY`
- `PUSHX_PSQL_QUERY_PARAMS`
- `PUSHX_PSQL_SSL_MODE`
- `PUSHX_PSQL_TLS_CERT`
- `PUSHX_PSQL_TLS_KEY`
- `PUSHX_PSQL_TLS_ROOT_CERT`
- `PUSHX_PSQL_USER`
- `PUSHX_PULSAR_ADDRESS`
- `PUSHX_PULSAR_AUTH_CERT_FILE`
- `PUSHX_PULSAR_AUTH_KEY_FILE`
- `PUSHX_PULSAR_AUTH_OAUTH_PARAMS`
- `PUSHX_PULSAR_AUTH_TOKEN`
- `PUSHX_PULSAR_AUTH_TOKEN_FILE`
- `PUSHX_PULSAR_PRODUCER_NAME`
- `PUSHX_PULSAR_TLS_ALLOW_INSECURE_CONNECTION`
- `PUSHX_PULSAR_TLS_TRUST_CERTS_FILE`
- `PUSHX_PULSAR_TLS_VALIDATE_HOSTNAME`
- `PUSHX_PULSAR_TOPIC`
- `PUSHX_RABBITMQ_QUEUE`
- `PUSHX_RABBITMQ_URL`
- `PUSHX_REDIS_ENABLE_TLS`
- `PUSHX_REDIS_HOST`
- `PUSHX_REDIS_KEY`
- `PUSHX_REDIS_MESSAGE_ID`
- `PUSHX_REDIS_PASSWORD`
- `PUSHX_REDIS_PORT`
- `PUSHX_REDIS_TLS_CA_FILE`
- `PUSHX_REDIS_TLS_CERT_FILE`
- `PUSHX_REDIS_TLS_INSECURE`
- `PUSHX_REDIS_TLS_KEY_FILE`

## Driver Examples

### ActiveMQ

The ActiveMQ driver will connect to the specified STOMP address and send the data to the specified queue / topic. TLS is optional and shown below, if not used the flags are not required.

```bash
echo hello | pushx \
    -driver activemq \
    -activemq-address localhost:61613 \
    -activemq-name my-queue \
    -activemq-enable-tls \
    -activemq-tls-ca-file /path/to/ca.pem \
    -activemq-tls-cert-file /path/to/cert.pem \
    -activemq-tls-key-file /path/to/key.pem
```

### AWS DynamoDB

The AWS DynamoDB driver will insert the specified JSON document into DynamoDB.

```bash
echo '{"hello": "world"}' | pushx \
    -driver aws-dynamo \
    -aws-dynamo-table my-table \
    -aws-region us-east-1 \
    -aws-role-arn arn:aws:iam::123456789012:role/my-role
```

### AWS S3

The S3 driver will upload the specified data to the specified S3 bucket.

```bash
pushx \
    -driver aws-s3 \
    -in-file /path/to/file.txt \
    -aws-s3-bucket my-bucket \
    -aws-s3-key 'example-object' \
    -aws-s3-acl public-read \
    -aws-s3-tags 'Hello=World,Foo=Bar'
```

### AWS SQS

The SQS driver will send the specified data to the specified SQS queue.

For cross-account access, you must provide the ARN of the role that has access to the queue, and the identity running pushx must be able to assume the target identity.

If running on a developer workstation, you will most likely want to pass your `~/.aws/config` identity. To do so, pass the `-aws-load-config` flag.

```bash
echo hello | pushx \
    -aws-sqs-queue-url https://sqs.us-east-1.amazonaws.com/123456789012/my-queue \
    -aws-role-arn arn:aws:iam::123456789012:role/my-role \
    -aws-region us-east-1 \
    -driver aws-sqs
```

### Cassandra

The Cassandra driver will submit the data to the specified keyspace table. If the `-cassandra-query` contains a `{{pushx_payload}}` placeholder, the entire input data will be substituted for the placeholder. However if the input data is a JSON document, value keys can be substituted using mustache-style syntax.

```bash
echo '{"hello": "world", "another": {"nested": "value"}}' | pushx \
    -cassandra-keyspace mykeyspace \
    -cassandra-consistency QUORUM \
    -cassandra-hosts "localhost:9042,another:9042" \
    -cassandra-query 'INSERT INTO mykeyspace.mytable (hello, another) VALUES (?, ?)' \
    -cassandra-params '{{hello}}, {{another.nested}}' \
    -driver cassandra
```

### Centauri

The `centauri` driver integrates with a [Centauri](https://centauri.sh) network to send the input data to the specified public key.

```bash
echo hello | pushx \
    -centauri-channel my-channel \
    -centauri-public-key "$(</path/to/public.pem)" \
    -centauri-peer-url https://api.test-peer1.centauri.sh \
    -driver centauri
```

### Cockroach

The Cockroach driver will insert the specified document into CockroachDB. If the input data is a JSON document, value keys can be substituted using mustache-style syntax.

```bash
echo '{"id": 1, "name": "hello", "another": "value"}' | pushx \
    -cockroach-host localhost \
    -cockroach-port 26257 \
    -cockroach-database mydb \
    -cockroach-user myuser \
    -cockroach-password mypassword \
    -cockroach-query 'INSERT INTO example (id, name, another) VALUES ($1, $2, $3)' \
    -cockroach-params "{{id}},{{name}},{{another}}" \
    -cockroach mysql
```

### Couchbase

The Couchbase driver will insert the specified document into the Couchbase bucket. Mustache syntax can be used to extract a value from the document to set as the document ID.

```bash
echo '{"hello": "world"}' | pushx \
    -couchbase-bucket my-bucket \
    -couchbase-user my-user \
    -couchbase-password my-password \
    -couchbase-scope my-scope \
    -couchbase-collection my-collection \
    -couchbase-address 'couchbase://localhost' \
    -couchbase-id '{{hello}}' \
    -driver couchbase
```

### Elasticsearch

The Elasticsearch driver will insert the specified document into Elasticsearch.

```bash
echo '{"hello": "world"}' | pushx \
    -elasticsearch-address https://localhost:9200 \
    -elasticsearch-username elastic \
    -elasticsearch-password elastic \
    -elasticsearch-tls-skip-verify \
    -elasticsearch-index my-index \
    -driver elasticsearch
```

### FS

The `fs` driver will write the input data to the specified file.

```bash
echo hello | pushx \
    -fs-folder /path/to/folder \
    -fs-key "my-file.txt" \
    -driver fs
```

### GCP BQ

The `gcp-bq` driver will insert the specified document into BigQuery. If the input data is a JSON document, value keys can be substituted using mustache-style syntax.

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials.json"
echo '{"id": 1, "name": "hello", "another": "value"}' | pushx \
    -gcp-project-id my-project \
    -gcp-bq-dataset my-dataset \
    -gcp-bq-table my-table \
    -gcp-bq-query "INSERT INTO mydatatest.mytable (id, name, another) VALUES ({{id}}, '{{name}}', '{{another}}')" \
    -driver gcp-bq
```

### GCP GCS

The GCS driver will upload the specified data to the specified GCS bucket.

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
pushx \
    -driver gcp-gcs \
    -gcp-project-id my-project \
    -payload-file my-payload.json \
    -gcp-gcs-bucket my-bucket \
    -gcp-gcs-key 'example.json'
```

### GCP Firestore

The GCP Firestore driver will insert the specified document into Firestore.

```bash
echo '{"id": 1, "name": "hello", "another": "value"}' | pushx \
    -driver gcp-firestore \
    -gcp-project-id my-project \
    -gcp-firestore-collection my-collection
```    

### GCP Pub/Sub

The GCP Pub/Sub driver will publish the specified data to the specified Pub/Sub topic.

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
echo hello | pushx \
    -gcp-project-id my-project \
    -gcp-pubsub-topic my-topic \
    -driver gcp-pubsub
```

### GitHub

The GitHub driver will insert the specified object `-github-file` into the provided GitHub repository `-github-repo`.  This can be done either on the same branch (`-github-ref` or `-github-base-branch`), or a new branch (`-github-branch`). If on a new branch, a pull request can be opened with `-github-open-pr`. If opening a new PR without a branch specified, a new branch name will be generated.

```bash
echo hello | pushx \
    -github-repo my-repo \
    -github-owner my-owner \
    -github-file path/to/my-file.txt \
    -github-token "$(</path/to/token.txt)" \
    -github-base-branch main \
    -github-branch my-branch \
    -github-open-pr \
    -driver github
```

### HTTP

The HTTP driver will connect to any HTTP(s) endpoint and submit the input data as a HTTP request. By default, the `POST` method is used. If using internal PKI, mTLS, or disabling TLS validation, pass the `-http-enable-tls` flag and the corresponding TLS flags.

```bash
echo hello | pushx \
    -http-url https://example.com/jobs \
    -http-headers 'ExampleToken:foobar,ExampleHeader:barfoo' \
    -driver http
```

### Kafka

The Kafka driver will submit the input data to the specified Kafka topic. Similar to the Pub/Sub drivers, if there are no messages in the topic when the process starts, it will wait for the first message. TLS and SASL authentication are optional.

```bash
echo hello | pushx \
    -kafka-brokers localhost:9092 \
    -kafka-topic my-topic \
    -kafka-enable-tls \
    -kafka-tls-ca-file /path/to/ca.pem \
    -kafka-tls-cert-file /path/to/cert.pem \
    -kafka-tls-key-file /path/to/key.pem \
    -kafka-enable-sasl \
    -kafka-sasl-type plain \
    -kafka-sasl-username my-username \
    -kafka-sasl-password my-password \
    -driver kafka
```

### MongoDB

The MongoDB driver will insert the specified document into MongoDB.

```bash
echo '{"id": 1, "name": "hello", "another": "value"}' | pushx \
    -mongo-collection my-collection \
    -mongo-database my-database \
    -mongo-host localhost \
    -mongo-port 27017 \
    -mongo-user my-user \
    -mongo-password my-password \
    -driver mongodb
```

### MSSQL

The MSSQL driver will insert the specified document into Microsoft SQL Server. If the input data is a JSON document, value keys can be substituted using mustache-style syntax.

```bash
echo '{"id": 1, "name": "hello", "another": "value"}' | pushx \
    -mssql-host localhost \
    -mssql-port 1433 \
    -mssql-database mydb \
    -mssql-user sa \
    -mssql-password 'mypassword!' \
    -mssql-query "INSERT INTO example (id, name, another) VALUES (?, ?, ?)" \
    -mssql-params "{{id}},{{name}},{{another}}" \
    -driver mssql
```

### MySQL

The MySQL driver will insert the specified document into MySQL. If the input data is a JSON document, value keys can be substituted using mustache-style syntax.

```bash
echo '{"id": 1, "name": "hello", "another": "value"}' | pushx \
    -mysql-host localhost \
    -mysql-port 3306 \
    -mysql-database mydb \
    -mysql-user myuser \
    -mysql-password mypassword \
    -mysql-query "INSERT INTO example (id, name, another) VALUES (?, ?, ?)" \
    -mysql-params "{{id}},{{name}},{{another}}" \
    -driver mysql
```

### NATS

The NATS driver will publish the specified data to the specified NATS topic.

```bash
echo hello | pushx \
    -nats-subject my-subject \
    -nats-url localhost:4222 \
    -nats-username my-user \
    -nats-password my-password \
    -nats-enable-tls \
    -nats-tls-ca-file /path/to/ca.pem \
    -nats-tls-cert-file /path/to/cert.pem \
    -nats-tls-key-file /path/to/key.pem \
    -driver nats
```

### NFS

The `nfs` driver will mount the specified NFS directory, and write the input data to the specified file.

```bash
echo hello | pushx \
    -nfs-host nfs.example.com \
    -nfs-target /path/to/nfs \
    -nfs-folder /path/to/folder/in/nfs \
    -nfs-key "my-object" \
    -driver nfs
```

### NSQ

The NSQ driver will connect to the specified `nsqlookupd` or `nsqd` endpoint and submit the input data as a `PUB` message.

```bash
echo hello | pushx \
    -nsq-nsqlookupd-address localhost:4161 \
    -nsq-topic my-topic \
    -driver nsq
```

### PostgreSQL

The PostgreSQL driver will insert the specified document into PostgreSQL. If the input data is a JSON document, value keys can be substituted using mustache-style syntax.

```bash
echo 'full payload' | pushx \
    -psql-host localhost \
    -psql-port 5432 \
    -psql-database mydb \
    -psql-user myuser \
    -psql-password mypassword \
    -psql-query "INSERT INTO example (payload) VALUES (?)" \
    -psql-params "{{pushx_payload}}" \
    -driver postgres
```

### Pulsar

The Pulsar driver will connect to the specified comma-separated Pulsar endpoint(s) and submit the input data as a `PUB` message.

```bash
echo hello | pushx \
    -pulsar-address localhost:6650,localhost:6651 \
    -pulsar-topic my-topic \
    -pulsar-producer-name my-producer \
    -pulsar-auth-cert-file /path/to/cert.pem \
    -pulsar-auth-key-file /path/to/key.pem \
    -pulsar-tls-trust-certs-file /path/to/trusted.pem \
    -driver pulsar
```

### RabbitMQ

The RabbitMQ driver will connect to the specified queue AMQP endpoint and submit the input data as a `PUB` message.

```bash
echo hello | pushx \
    -rabbitmq-url amqp://guest:guest@localhost:5672 \
    -rabbitmq-queue my-queue \
    -driver rabbitmq
```

### Redis List

The Redis List driver will connect to the specified Redis server and push the input data to the specified list.

```bash
echo hello | pushx \
    -redis-host localhost \
    -redis-port 6379 \
    -redis-key my-list \
    -driver redis-list
```

### Redis Pub/Sub

The Redis Pub/Sub driver will connect to the specified Redis server and publish the input data to the specified topic.

```bash
echo hello | pushx \
    -redis-host localhost \
    -redis-port 6379 \
    -redis-key my-subscription \
    -driver redis-pubsub
```

### Redis Stream

The Redis Stream driver will connect to the specified Redis server and push the input data to the specified stream.

```bash
echo '{"id": 1, "name": "hello", "another": "value"}' | pushx \
    -redis-host localhost \
    -redis-port 6379 \
    -redis-key my-stream \
    -driver redis-stream
```

### Local

The local driver is a simple wrapper around the process to execute, primarily for local testing. It does not communicate with any queue, and simply writes the input data to the specified output (default is stdout).

```bash
echo hello | pushx \
    -driver local
```
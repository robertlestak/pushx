package flags

var (
	GCPProjectID = FlagSet.String("gcp-project-id", "", "GCP project ID")
	GCPTopic     = FlagSet.String("gcp-pubsub-topic", "", "GCP Pub/Sub topic name")

	GCPFirestoreCollection = FlagSet.String("gcp-firestore-collection", "", "GCP Firestore collection")
	GCPFirestoreID         = FlagSet.String("gcp-firestore-id", "", "GCP Firestore document ID. If empty, a new document ID will be created")

	GCPGCSBucket = FlagSet.String("gcp-gcs-bucket", "", "GCP GCS bucket")
	GCPGCSKey    = FlagSet.String("gcp-gcs-key", "", "GCP GCS key")

	GCPBQQuery = FlagSet.String("gcp-bq-query", "", "GCP BigQuery query")
)

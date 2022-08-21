package flags

var (
	AWSRegion     = FlagSet.String("aws-region", "", "AWS region")
	AWSLoadConfig = FlagSet.Bool("aws-load-config", false, "load AWS config from ~/.aws/config")
	AWSRoleARN    = FlagSet.String("aws-role-arn", "", "AWS role ARN")
	SQSQueueURL   = FlagSet.String("aws-sqs-queue-url", "", "AWS SQS queue URL")

	AWSDynamoTable = FlagSet.String("aws-dynamo-table", "", "AWS DynamoDB table name")

	AWSS3Bucket = FlagSet.String("aws-s3-bucket", "", "AWS S3 bucket")
	AWSS3Key    = FlagSet.String("aws-s3-key", "", "AWS S3 key")
	AWSS3ACL    = FlagSet.String("aws-s3-acl", "", "AWS S3 ACL")
	AWSS3Tags   = FlagSet.String("aws-s3-tags", "", "AWS S3 tags. Comma separated list of key=value pairs")
)

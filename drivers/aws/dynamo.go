package aws

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/uuid"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type Dynamo struct {
	Client  *dynamodb.DynamoDB
	sts     *STSSession
	Table   string
	Region  string
	RoleARN string
}

func (d *Dynamo) LogIdentity() error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "LogIdentity",
	})
	l.Debug("LogIdentity")
	streq := &sts.GetCallerIdentityInput{}
	var sc *sts.STS
	if d.sts.Config != nil {
		sc = sts.New(d.sts.Session, d.sts.Config)
	} else {
		sc = sts.New(d.sts.Session)
	}
	r, err := sc.GetCallerIdentity(streq)
	if err != nil {
		l.Errorf("%+v", err)
	} else {
		l.Debugf("%+v", r)
	}
	return nil
}

func (d *Dynamo) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"fn":  "LoadEnv",
		"pkg": "aws",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"AWS_REGION") != "" {
		d.Region = os.Getenv(prefix + "AWS_REGION")
	}
	if os.Getenv(prefix+"AWS_ROLE_ARN") != "" {
		d.RoleARN = os.Getenv(prefix + "AWS_SQS_ROLE_ARN")
	}
	if os.Getenv(prefix+"AWS_DYNAMO_TABLE") != "" {
		d.Table = os.Getenv(prefix + "AWS_DYNAMO_TABLE")
	}
	if os.Getenv(prefix+"AWS_LOAD_CONFIG") != "" || os.Getenv("AWS_SDK_LOAD_CONFIG") != "" {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	}
	return nil
}

func (d *Dynamo) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"fn":  "LoadFlags",
		"pkg": "aws",
	})
	l.Debug("LoadFlags")
	d.Table = *flags.AWSDynamoTable
	d.Region = *flags.AWSRegion
	d.RoleARN = *flags.AWSRoleARN
	if flags.AWSLoadConfig != nil && *flags.AWSLoadConfig {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	}
	return nil
}

func (d *Dynamo) Init() error {
	l := log.WithFields(
		log.Fields{
			"fn":  "CreateAWSSession",
			"pkg": "aws",
		},
	)
	l.Debug("CreateAWSSession")
	if d.Region == "" {
		d.Region = os.Getenv("AWS_REGION")
	}
	if d.Region == "" {
		d.Region = "us-east-1"
	}
	cfg := &aws.Config{
		Region: aws.String(d.Region),
	}
	sess, err := session.NewSession(cfg)
	reqId := uuid.New().String()
	if d.RoleARN != "" {
		l.Debug("CreateAWSSession roleArn=%s requestId=%s", d.RoleARN, reqId)
		creds := stscreds.NewCredentials(sess, d.RoleARN, func(p *stscreds.AssumeRoleProvider) {
			p.RoleSessionName = "pushx-" + reqId
		})
		cfg.Credentials = creds
	}
	if err != nil {
		l.Errorf("%+v", err)
		if err := d.LogIdentity(); err != nil {
			l.Errorf("%+v", err)
		}
		return err

	}
	d.Client = dynamodb.New(sess, cfg)
	d.sts = &STSSession{
		Session: sess,
		Config:  cfg,
	}
	return err
}

func (d *Dynamo) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"fn":  "Push",
		"pkg": "aws",
	})
	l.Debug("Push")
	var err error
	var data []byte
	if data, err = ioutil.ReadAll(r); err != nil {
		l.Errorf("%+v", err)
		return err
	}
	var m map[string]interface{}
	if err = json.Unmarshal(data, &m); err != nil {
		l.Errorf("%+v", err)
		return err
	}
	av, err := dynamodbattribute.MarshalMap(m)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	if _, err = d.Client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.Table),
	}); err != nil {
		l.Errorf("%+v", err)
		return err
	}
	return nil
}

func (d *Dynamo) Cleanup() error {
	l := log.WithFields(log.Fields{
		"fn":  "Cleanup",
		"pkg": "aws",
	})
	l.Debug("Cleanup")
	return nil
}

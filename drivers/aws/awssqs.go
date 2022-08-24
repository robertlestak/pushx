package aws

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/uuid"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type STSSession struct {
	Session *session.Session
	Config  *aws.Config
}

type SQS struct {
	Client  *sqs.SQS
	sts     *STSSession
	Queue   string
	Region  string
	RoleARN string
}

func (d *SQS) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"AWS_REGION") != "" {
		d.Region = os.Getenv(prefix + "AWS_REGION")
	}
	if os.Getenv(prefix+"AWS_ROLE_ARN") != "" {
		d.RoleARN = os.Getenv(prefix + "AWS_ROLE_ARN")
	}
	if os.Getenv(prefix+"AWS_SQS_QUEUE_URL") != "" {
		d.Queue = os.Getenv(prefix + "AWS_SQS_QUEUE_URL")
	}
	if os.Getenv(prefix+"AWS_LOAD_CONFIG") != "" || os.Getenv("AWS_SDK_LOAD_CONFIG") != "" {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	}
	return nil
}

func (d *SQS) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.Queue = *flags.SQSQueueURL
	d.Region = *flags.AWSRegion
	d.RoleARN = *flags.AWSRoleARN
	if flags.AWSLoadConfig != nil && *flags.AWSLoadConfig {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	}
	return nil
}

func (d *SQS) LogIdentity() error {
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

func (d *SQS) Init() error {
	l := log.WithFields(
		log.Fields{
			"pkg": "aws",
			"fn":  "CreateAWSSession",
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
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
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
	d.Client = sqs.New(sess, cfg)
	d.sts = &STSSession{
		Session: sess,
		Config:  cfg,
	}
	return err
}

func (d *SQS) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "Push",
	})
	l.Debug("Push")
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	req := &sqs.SendMessageInput{
		MessageBody: aws.String(strings.TrimSpace(string(bd))),
		QueueUrl:    aws.String(d.Queue),
	}
	_, err = d.Client.SendMessage(req)
	if err != nil {
		l.Errorf("%+v", err)
	}
	return err
}

func (d *SQS) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "Cleanup",
	})
	l.Debug("Cleanup")
	return nil
}

package aws

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/uuid"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type S3 struct {
	Client  *s3manager.Uploader
	sts     *STSSession
	Bucket  string
	Key     string
	Region  string
	RoleARN string
	ACL     string
	Tags    map[string]string
}

func (d *S3) LogIdentity() error {
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

func (d *S3) LoadEnv(prefix string) error {
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
	if os.Getenv(prefix+"AWS_S3_BUCKET") != "" {
		d.Bucket = os.Getenv(prefix + "AWS_S3_BUCKET")
	}
	if os.Getenv(prefix+"AWS_S3_KEY") != "" {
		d.Key = os.Getenv(prefix + "AWS_S3_KEY")
	}
	if os.Getenv(prefix+"AWS_S3_ACL") != "" {
		d.ACL = os.Getenv(prefix + "AWS_S3_ACL")
	}
	if os.Getenv(prefix+"AWS_S3_TAGS") != "" {
		d.Tags = make(map[string]string)
		for _, tag := range strings.Split(os.Getenv(prefix+"AWS_S3_TAGS"), ",") {
			parts := strings.Split(tag, "=")
			if len(parts) == 2 {
				d.Tags[parts[0]] = parts[1]
			}
		}
	}
	if os.Getenv(prefix+"AWS_LOAD_CONFIG") != "" || os.Getenv("AWS_SDK_LOAD_CONFIG") != "" {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	}
	return nil
}

func (d *S3) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.Bucket = *flags.AWSS3Bucket
	d.Key = *flags.AWSS3Key
	d.Region = *flags.AWSRegion
	d.RoleARN = *flags.AWSRoleARN
	d.ACL = *flags.AWSS3ACL
	d.Tags = make(map[string]string)
	for _, tag := range strings.Split(*flags.AWSS3Tags, ",") {
		parts := strings.Split(tag, "=")
		if len(parts) == 2 {
			d.Tags[parts[0]] = parts[1]
		}
	}
	if flags.AWSLoadConfig != nil && *flags.AWSLoadConfig {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	}
	return nil
}

func (d *S3) Init() error {
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
	reqId := uuid.New().String()
	if d.RoleARN != "" {
		l.Debug("CreateAWSSession roleArn=%s requestId=%s", d.RoleARN, reqId)
		creds := stscreds.NewCredentials(sess, d.RoleARN, func(p *stscreds.AssumeRoleProvider) {
			p.RoleSessionName = "pushx-" + reqId
		})
		cfg.Credentials = creds
		sess.Config = cfg
	}
	if err != nil {
		l.Errorf("%+v", err)
		if err := d.LogIdentity(); err != nil {
			l.Errorf("%+v", err)
		}
		return err

	}
	d.Client = s3manager.NewUploader(sess)
	d.sts = &STSSession{
		Session: sess,
		Config:  cfg,
	}
	return err
}

type reader struct {
	r io.Reader
}

func (r *reader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (d *S3) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "Push",
	})
	l.Debug("Push")
	if d.Bucket == "" {
		return errors.New("bucket not set")
	}
	if d.Key == "" {
		return errors.New("key not set")
	}
	req := &s3manager.UploadInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.Key),
		Body:   &reader{r},
	}
	if d.ACL != "" {
		req.ACL = aws.String(d.ACL)
	}
	if d.Tags != nil {
		// encode tags as url query string
		var buf bytes.Buffer
		for k, v := range d.Tags {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
		req.Tagging = aws.String(buf.String())
	}
	_, err := d.Client.Upload(req)
	if err != nil {
		l.Errorf("%+v", err)
		if err := d.LogIdentity(); err != nil {
			l.Errorf("%+v", err)
		}
		return err
	}
	return nil
}

func (d *S3) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "aws",
		"fn":  "Cleanup",
	})
	l.Debug("Cleanup")
	return nil
}

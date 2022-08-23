package couchbase

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/couchbase/gocb/v2"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/schema"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Couchbase struct {
	Client     *gocb.Cluster
	Address    string
	User       *string
	Password   *string
	BucketName *string
	Scope      *string
	Collection *string
	ID         *string
	// TLS
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCert     *string
	TLSKey      *string
	TLSCA       *string
}

func (d *Couchbase) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "couchbase",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"COUCHBASE_USER") != "" {
		v := os.Getenv(prefix + "COUCHBASE_USER")
		d.User = &v
	}
	if os.Getenv(prefix+"COUCHBASE_PASSWORD") != "" {
		v := os.Getenv(prefix + "COUCHBASE_PASSWORD")
		d.Password = &v
	}
	if os.Getenv(prefix+"COUCHBASE_COLLECTION") != "" {
		c := os.Getenv(prefix + "COUCHBASE_COLLECTION")
		d.Collection = &c
	}
	if os.Getenv(prefix+"COUCHBASE_SCOPE") != "" {
		s := os.Getenv(prefix + "COUCHBASE_SCOPE")
		d.Scope = &s
	}
	if os.Getenv(prefix+"COUCHBASE_BUCKET_NAME") != "" {
		b := os.Getenv(prefix + "COUCHBASE_BUCKET_NAME")
		d.BucketName = &b
	}
	if os.Getenv(prefix+"COUCHBASE_ID") != "" {
		i := os.Getenv(prefix + "COUCHBASE_ID")
		d.ID = &i
	}
	if os.Getenv(prefix+"COUCHBASE_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"COUCHBASE_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"COUCHBASE_TLS_INSECURE") != "" {
		v := os.Getenv(prefix+"COUCHBASE_TLS_INSECURE") == "true"
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"COUCHBASE_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "COUCHBASE_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"COUCHBASE_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "COUCHBASE_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	if os.Getenv(prefix+"COUCHBASE_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "COUCHBASE_TLS_CA_FILE")
		d.TLSCA = &v
	}
	return nil
}

func (d *Couchbase) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "couchbase",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.User = flags.CouchbaseUser
	d.Password = flags.CouchbasePassword
	d.BucketName = flags.CouchbaseBucketName
	d.Scope = flags.CouchbaseScope
	d.Address = *flags.CouchbaseAddress
	d.Collection = flags.CouchbaseCollection
	d.ID = flags.CouchbaseID
	d.EnableTLS = flags.CouchbaseEnableTLS
	d.TLSInsecure = flags.CouchbaseTLSInsecure
	d.TLSCert = flags.CouchbaseCertFile
	d.TLSKey = flags.CouchbaseKeyFile
	d.TLSCA = flags.CouchbaseCAFile
	return nil
}

func (d *Couchbase) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "couchbase",
		"fn":  "Init",
	})
	l.Debug("Initializing couchbase client")
	opts := gocb.ClusterOptions{}
	if d.User != nil && *d.User != "" && d.Password != nil && *d.Password != "" {
		opts.Authenticator = gocb.PasswordAuthenticator{
			Username: *d.User,
			Password: *d.Password,
		}
	}
	if d.EnableTLS != nil && *d.EnableTLS {
		tc, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
		if err != nil {
			return err
		}
		sc := gocb.SecurityConfig{
			TLSSkipVerify: d.TLSInsecure != nil && *d.TLSInsecure,
			TLSRootCAs:    tc.RootCAs,
		}
		opts.SecurityConfig = sc
	}
	cluster, err := gocb.Connect(d.Address, opts)
	if err != nil {
		l.Error(err)
		return err
	}
	d.Client = cluster
	return nil
}

func getID(data map[string]any, in string) string {
	l := log.WithFields(log.Fields{
		"pkg": "couchbase",
		"fn":  "getID",
	})
	l.Debug("Getting ID")
	if in == "" {
		return in
	}
	if !strings.Contains(in, "{{") {
		return in
	}
	jd, err := json.Marshal(data)
	if err != nil {
		l.Error(err)
		return ""
	}
	return schema.ReplaceParamsString(jd, in)
}

func (d *Couchbase) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "couchbase",
		"fn":  "Push",
	})
	l.Debug("Inserting data")
	var message map[string]interface{}
	err := json.NewDecoder(r).Decode(&message)
	if err != nil {
		l.Error(err)
		return err
	}
	if d.BucketName == nil || *d.BucketName == "" {
		l.Error("bucket name is not set")
		return errors.New("bucket name is not set")
	}
	if d.ID == nil || *d.ID == "" {
		l.Error("id is not set")
		return errors.New("id is not set")
	}
	i := getID(message, *d.ID)
	if i == "" {
		l.Error("id is empty")
		return errors.New("id is empty")
	}
	d.ID = &i
	bucket := d.Client.Bucket(*d.BucketName)
	coll := bucket.Scope(*d.Scope).Collection(*d.Collection)
	_, err = coll.Upsert(*d.ID, message, nil)
	if err != nil {
		l.Error(err)
		return err
	}
	return nil
}

func (d *Couchbase) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "couchbase",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up couchbase")
	if d.Client == nil {
		return nil
	}
	d.Client.Close(nil)
	l.Debug("Cleaned up couchbase")
	return nil
}

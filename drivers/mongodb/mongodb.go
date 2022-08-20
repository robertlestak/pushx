package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client     *mongo.Client
	Host       string
	Port       int
	User       string
	Password   string
	DB         string
	Collection string
	AuthSource string
	// TLS
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCert     *string
	TLSKey      *string
	TLSCA       *string
}

func (d *Mongo) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "mongo",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"MONGO_HOST") != "" {
		d.Host = os.Getenv(prefix + "MONGO_HOST")
	}
	if os.Getenv(prefix+"MONGO_PORT") != "" {
		pval, err := strconv.Atoi(os.Getenv(prefix + "MONGO_PORT"))
		if err != nil {
			return err
		}
		d.Port = pval
	}
	if os.Getenv(prefix+"MONGO_USER") != "" {
		d.User = os.Getenv(prefix + "MONGO_USER")
	}
	if os.Getenv(prefix+"MONGO_PASSWORD") != "" {
		d.Password = os.Getenv(prefix + "MONGO_PASSWORD")
	}
	if os.Getenv(prefix+"MONGO_DATABASE") != "" {
		d.DB = os.Getenv(prefix + "MONGO_DATABASE")
	}
	if os.Getenv(prefix+"MONGO_COLLECTION") != "" {
		c := os.Getenv(prefix + "MONGO_COLLECTION")
		d.Collection = strings.TrimSpace(c)
	}
	if os.Getenv(prefix+"MONGO_AUTH_SOURCE") != "" {
		d.AuthSource = os.Getenv(prefix + "MONGO_AUTH_SOURCE")
	}
	if os.Getenv(prefix+"MONGO_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"MONGO_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"MONGO_TLS_INSECURE") != "" {
		v := os.Getenv(prefix+"MONGO_TLS_INSECURE") == "true"
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"MONGO_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "MONGO_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"MONGO_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "MONGO_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	if os.Getenv(prefix+"MONGO_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "MONGO_TLS_CA_FILE")
		d.TLSCA = &v
	}
	return nil
}

func (d *Mongo) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "mongo",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	pv, err := strconv.Atoi(*flags.MongoPort)
	if err != nil {
		return err
	}
	d.Host = *flags.MongoHost
	d.Port = pv
	d.User = *flags.MongoUser
	d.Password = *flags.MongoPassword
	d.AuthSource = *flags.MongoAuthSource
	d.DB = *flags.MongoDatabase
	d.Collection = *flags.MongoCollection
	d.EnableTLS = flags.MongoEnableTLS
	d.TLSInsecure = flags.MongoTLSInsecure
	d.TLSCert = flags.MongoCertFile
	d.TLSKey = flags.MongoKeyFile
	d.TLSCA = flags.MongoCAFile
	return nil
}

func (d *Mongo) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "mongo",
		"fn":  "Init",
	})
	l.Debug("Initializing mongo client")
	var err error
	var uri string
	if d.AuthSource != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s", d.User, d.Password, d.Host, d.Port, d.DB, d.AuthSource)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.DB)
	}
	l.Debug("uri: ", uri)
	opts := options.Client().ApplyURI(uri)
	if d.EnableTLS != nil && *d.EnableTLS {
		l.Debug("TLS enabled")
		tc, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
		if err != nil {
			return err
		}
		opts.SetTLSConfig(tc)
	}
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		l.Error(err)
		return err
	}
	d.Client = client
	// ping the database to check if it is alive
	err = d.Client.Ping(context.TODO(), nil)
	if err != nil {
		l.Error(err)
		return err
	}
	return nil
}

func (d *Mongo) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "mongo",
		"fn":  "Push",
	})
	l.Debug("Inserting data")
	var message map[string]interface{}
	err := json.NewDecoder(r).Decode(&message)
	if err != nil {
		l.Error(err)
		return err
	}
	collection := d.Client.Database(d.DB).Collection(d.Collection)
	res, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Inserted 1 document: ", res.InsertedID)
	return nil
}

func (d *Mongo) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "mongo",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up mongo")
	if d.Client == nil {
		return nil
	}
	err := d.Client.Disconnect(context.TODO())
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Cleaned up mongo")
	return nil
}

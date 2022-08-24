package elasticsearch

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Elasticsearch struct {
	Client   *elasticsearch8.Client
	Address  string
	Username string
	Password string
	// TLS
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCert     *string
	TLSKey      *string
	TLSCA       *string
	Index       *string
	Key         *string
}

func (d *Elasticsearch) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "elasticsearch",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	if os.Getenv(prefix+"ELASTICSEARCH_ADDRESS") != "" {
		d.Address = os.Getenv(prefix + "ELASTICSEARCH_ADDRESS")
	}
	if os.Getenv(prefix+"ELASTICSEARCH_USERNAME") != "" {
		d.Username = os.Getenv(prefix + "ELASTICSEARCH_USERNAME")
	}
	if os.Getenv(prefix+"ELASTICSEARCH_PASSWORD") != "" {
		d.Password = os.Getenv(prefix + "ELASTICSEARCH_PASSWORD")
	}
	if os.Getenv(prefix+"ELASTICSEARCH_TLS_SKIP_VERIFY") != "" {
		v := os.Getenv(prefix+"ELASTICSEARCH_TLS_SKIP_VERIFY") == "true"
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"ELASTICSEARCH_INDEX") != "" {
		v := os.Getenv(prefix + "ELASTICSEARCH_INDEX")
		d.Index = &v
	}
	if os.Getenv(prefix+"ELASTICSEARCH_DOC_ID") != "" {
		v := os.Getenv(prefix + "ELASTICSEARCH_DOC_ID")
		d.Key = &v
	}
	if os.Getenv(prefix+"ELASTICSEARCH_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"ELASTICSEARCH_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"ELASTICSEARCH_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "ELASTICSEARCH_TLS_CA_FILE")
		d.TLSCA = &v
	}
	if os.Getenv(prefix+"ELASTICSEARCH_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "ELASTICSEARCH_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"ELASTICSEARCH_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "ELASTICSEARCH_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	return nil
}

func (d *Elasticsearch) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "elasticsearch",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.Address = *flags.ElasticsearchAddress
	d.Username = *flags.ElasticsearchUsername
	d.Password = *flags.ElasticsearchPassword
	d.Key = flags.ElasticsearchDocID
	d.TLSInsecure = flags.ElasticsearchTLSSkipVerify
	d.EnableTLS = flags.ElasticsearchEnableTLS
	d.TLSCert = flags.ElasticsearchCertFile
	d.TLSKey = flags.ElasticsearchKeyFile
	d.TLSCA = flags.ElasticsearchCAFile
	d.Index = flags.ElasticsearchIndex
	return nil
}

func (d *Elasticsearch) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "elasticsearch",
		"fn":  "Init",
	})
	l.Debug("Initializing elasticsearch driver")
	tc, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
	if err != nil {
		return err
	}
	client, err := elasticsearch8.NewClient(elasticsearch8.Config{
		Transport: &http.Transport{
			TLSClientConfig: tc,
		},
		Addresses: []string{d.Address},
		Username:  d.Username,
		Password:  d.Password,
	})
	if err != nil {
		l.Errorf("error creating client: %v", err)
		return err
	}
	d.Client = client
	return nil
}

func (d *Elasticsearch) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "elasticsearch",
		"fn":  "Push",
	})
	l.Debug("Putting work")
	put := esapi.IndexRequest{
		DocumentID: *d.Key,
	}
	if d.Index != nil {
		put.Index = *d.Index
	}
	put.Body = r
	putResponse, err := put.Do(context.Background(), d.Client)
	if err != nil {
		l.Errorf("error putting work: %v", err)
		return err
	}
	if putResponse.StatusCode != 200 && putResponse.StatusCode != 201 {
		bd, err := ioutil.ReadAll(putResponse.Body)
		if err != nil {
			l.Errorf("error reading response body: %v", err)
			return err
		}
		l.Errorf("error putting work(%d): %v", putResponse.StatusCode, string(bd))
		return errors.New("error putting work")
	}
	return nil
}

func (d *Elasticsearch) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "elasticsearch",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up")
	l.Debug("Cleaned up")
	return nil
}

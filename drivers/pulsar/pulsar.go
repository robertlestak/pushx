package pulsar

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	pulsarlog "github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type Pulsar struct {
	Client          pulsar.Client
	Address         string
	Topic           *string
	Name            *string
	AuthToken       *string
	AuthTokenFile   *string
	AuthCertPath    *string
	AuthKeyPath     *string
	AuthOAuthParams *map[string]string
	// TLS
	TLSTrustCertsFilePath      *string
	TLSAllowInsecureConnection *bool
	TLSValidateHostname        *bool
}

func (d *Pulsar) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "pulsar",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	if os.Getenv(prefix+"PULSAR_ADDRESS") != "" {
		d.Address = os.Getenv(prefix + "PULSAR_ADDRESS")
	}
	if os.Getenv(prefix+"PULSAR_TOPIC") != "" {
		v := os.Getenv(prefix + "PULSAR_TOPIC")
		d.Topic = &v
	}
	if os.Getenv(prefix+"PULSAR_PRODUCER_NAME") != "" {
		v := os.Getenv(prefix + "PULSAR_PRODUCER_NAME")
		d.Name = &v
	}
	if os.Getenv(prefix+"PULSAR_TLS_TRUST_CERTS_FILE") != "" {
		v := os.Getenv(prefix + "PULSAR_TLS_TRUST_CERTS_FILE")
		d.TLSTrustCertsFilePath = &v
	}
	if os.Getenv(prefix+"PULSAR_TLS_ALLOW_INSECURE_CONNECTION") != "" {
		v := os.Getenv(prefix+"PULSAR_TLS_ALLOW_INSECURE_CONNECTION") == "true"
		d.TLSAllowInsecureConnection = &v
	}
	if os.Getenv(prefix+"PULSAR_TLS_VALIDATE_HOSTNAME") != "" {
		v := os.Getenv(prefix+"PULSAR_TLS_VALIDATE_HOSTNAME") == "true"
		d.TLSValidateHostname = &v
	}
	if os.Getenv(prefix+"PULSAR_AUTH_TOKEN") != "" {
		v := os.Getenv(prefix + "PULSAR_AUTH_TOKEN")
		d.AuthToken = &v
	}
	if os.Getenv(prefix+"PULSAR_AUTH_TOKEN_FILE") != "" {
		v := os.Getenv(prefix + "PULSAR_AUTH_TOKEN_FILE")
		d.AuthTokenFile = &v
	}
	if os.Getenv(prefix+"PULSAR_AUTH_CERT_FILE") != "" {
		v := os.Getenv(prefix + "PULSAR_AUTH_CERT_FILE")
		d.AuthCertPath = &v
	}
	if os.Getenv(prefix+"PULSAR_AUTH_KEY_FILE") != "" {
		v := os.Getenv(prefix + "PULSAR_AUTH_KEY_FILE")
		d.AuthKeyPath = &v
	}
	if os.Getenv(prefix+"PULSAR_AUTH_OAUTH_PARAMS") != "" {
		v := os.Getenv(prefix + "PULSAR_AUTH_OAUTH_PARAMS")
		var m map[string]string
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return err
		}
		d.AuthOAuthParams = &m
	}
	l.Debug("Loaded environment")
	return nil
}

func (d *Pulsar) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "pulsar",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.Address = *flags.PulsarAddress
	d.Topic = flags.PulsarTopic
	d.Name = flags.PulsarProducerName
	d.TLSTrustCertsFilePath = flags.PulsarTLSTrustCertsFilePath
	d.TLSAllowInsecureConnection = flags.PulsarTLSAllowInsecureConnection
	d.TLSValidateHostname = flags.PulsarTLSValidateHostname
	d.AuthToken = flags.PulsarAuthToken
	d.AuthTokenFile = flags.PulsarAuthTokenFile
	d.AuthCertPath = flags.PulsarAuthCertFile
	d.AuthKeyPath = flags.PulsarAuthKeyFile
	oauthParams := make(map[string]string)
	if flags.PulsarAuthOAuthParams != nil && *flags.PulsarAuthOAuthParams != "" {
		if err := json.Unmarshal([]byte(*flags.PulsarAuthOAuthParams), &oauthParams); err != nil {
			l.WithError(err).Error("Failed to parse PulsarAuthOAuthParams")
			return err
		}
		d.AuthOAuthParams = &oauthParams
	}
	l.Debug("Loaded flags")
	return nil
}

func (d *Pulsar) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "pulsar",
		"fn":  "Init",
	})

	l.Debug("Initializing pulsar driver")
	ll := log.StandardLogger()
	ll.Level = log.FatalLevel
	lr := pulsarlog.NewLoggerWithLogrus(ll)
	up := "pulsar://"
	opts := pulsar.ClientOptions{
		OperationTimeout:           30 * time.Second,
		ConnectionTimeout:          30 * time.Second,
		TLSTrustCertsFilePath:      *d.TLSTrustCertsFilePath,
		TLSAllowInsecureConnection: *d.TLSAllowInsecureConnection,
		TLSValidateHostname:        *d.TLSValidateHostname,
		Logger:                     lr,
	}
	if *d.AuthCertPath != "" && *d.AuthKeyPath != "" {
		opts.Authentication = pulsar.NewAuthenticationTLS(*d.AuthCertPath, *d.AuthKeyPath)
		up = "pulsar+ssl://"
	}
	if *d.AuthToken != "" {
		opts.Authentication = pulsar.NewAuthenticationToken(*d.AuthToken)
	}
	if *d.AuthTokenFile != "" {
		opts.Authentication = pulsar.NewAuthenticationTokenFromFile(*d.AuthTokenFile)
	}
	if *d.AuthOAuthParams != nil && len(*d.AuthOAuthParams) > 0 {
		opts.Authentication = pulsar.NewAuthenticationOAuth2(*d.AuthOAuthParams)
	}
	opts.URL = up + d.Address
	client, err := pulsar.NewClient(opts)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	d.Client = client
	l.Debug("Initialized pulsar driver")
	return nil
}

func (d *Pulsar) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "pulsar",
		"fn":  "Push",
	})
	l.Debug("Pushing data to pulsar")
	topic := *d.Topic
	opts := pulsar.ProducerOptions{
		Topic: topic,
	}
	producer, err := d.Client.CreateProducer(opts)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: bd,
	})
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	l.Debug("Pushed data to pulsar")
	return nil
}

func (d *Pulsar) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "pulsar",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up")
	d.Client.Close()
	l.Debug("Cleaned up")
	return nil
}

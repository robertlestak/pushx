package nats

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type NATS struct {
	Client      *nats.Conn
	URL         string
	Subject     *string
	CredsFile   *string
	JWTFile     *string
	NKeyFile    *string
	Username    *string
	Password    *string
	Token       *string
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCA       *string
	TLSCert     *string
	TLSKey      *string
}

func (d *NATS) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "nats",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	if os.Getenv(prefix+"NATS_URL") != "" {
		d.URL = os.Getenv(prefix + "NATS_URL")
	}
	if os.Getenv(prefix+"NATS_SUBJECT") != "" {
		v := os.Getenv(prefix + "NATS_SUBJECT")
		d.Subject = &v
	}
	if os.Getenv(prefix+"NATS_CREDS_FILE") != "" {
		v := os.Getenv(prefix + "NATS_CREDS_FILE")
		d.CredsFile = &v
	}
	if os.Getenv(prefix+"NATS_JWT_FILE") != "" {
		v := os.Getenv(prefix + "NATS_JWT_FILE")
		d.JWTFile = &v
	}
	if os.Getenv(prefix+"NATS_NKEY_FILE") != "" {
		v := os.Getenv(prefix + "NATS_NKEY_FILE")
		d.NKeyFile = &v
	}
	if os.Getenv(prefix+"NATS_USERNAME") != "" {
		v := os.Getenv(prefix + "NATS_USERNAME")
		d.Username = &v
	}
	if os.Getenv(prefix+"NATS_PASSWORD") != "" {
		v := os.Getenv(prefix + "NATS_PASSWORD")
		d.Password = &v
	}
	if os.Getenv(prefix+"NATS_TOKEN") != "" {
		v := os.Getenv(prefix + "NATS_TOKEN")
		d.Token = &v
	}
	if os.Getenv(prefix+"NATS_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"NATS_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"NATS_TLS_INSECURE") != "" {
		v := os.Getenv(prefix+"NATS_TLS_INSECURE") == "true"
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"NATS_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "NATS_TLS_CA_FILE")
		d.TLSCA = &v
	}
	if os.Getenv(prefix+"NATS_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "NATS_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"NATS_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "NATS_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	return nil
}

func (d *NATS) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "nats",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.URL = *flags.NATSURL
	d.Subject = flags.NATSSubject
	d.CredsFile = flags.NATSCredsFile
	d.JWTFile = flags.NATSJWTFile
	d.NKeyFile = flags.NATSNKeyFile
	d.Username = flags.NATSUsername
	d.Password = flags.NATSPassword
	d.Token = flags.NATSToken
	d.EnableTLS = flags.NATSEnableTLS
	d.TLSInsecure = flags.NATSTLSInsecure
	d.TLSCA = flags.NATSTLSCAFile
	d.TLSCert = flags.NATSTLSCertFile
	d.TLSKey = flags.NATSTLSKeyFile
	l.Debug("Loaded flags")
	return nil
}

func (d *NATS) authOpts() []nats.Option {
	l := log.WithFields(log.Fields{
		"pkg": "nats",
		"fn":  "authOpts",
	})
	l.Debug("Creating auth options")
	opts := []nats.Option{}
	if d.CredsFile != nil && *d.CredsFile != "" {
		l.Debug("Enabling creds file")
		opts = append(opts, nats.UserCredentials(*d.CredsFile))
	}
	if d.Username != nil && *d.Username != "" {
		l.Debug("Enabling username")
		opts = append(opts, nats.UserInfo(*d.Username, *d.Password))
	}
	if d.Token != nil && *d.Token != "" {
		l.Debug("Enabling token")
		opts = append(opts, nats.Token(*d.Token))
	}
	if d.JWTFile != nil && *d.JWTFile != "" && d.NKeyFile != nil && *d.NKeyFile != "" {
		l.Debug("Enabling JWT file")
		opts = append(opts, nats.UserCredentials(*d.JWTFile, *d.NKeyFile))
	}
	l.Debug("Created auth options")
	return opts
}

func (d *NATS) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "nats",
		"fn":  "Init",
	})
	l.Debug("Initializing nats driver")
	if d.URL == "" {
		l.Error("url is empty")
		return errors.New("url is empty")
	}
	opts := []nats.Option{}
	if d.EnableTLS != nil && *d.EnableTLS {
		l.Debug("Enabling TLS")
		tc, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
		if err != nil {
			l.Errorf("%+v", err)
			return err
		}
		opts = append(opts, nats.Secure(tc))
	}
	opts = append(opts, d.authOpts()...)
	nc, err := nats.Connect(d.URL, opts...)
	if err != nil {
		l.Errorf("error connecting to nats: %v", err)
		return err
	}
	d.Client = nc
	return nil
}

func (d *NATS) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "nats",
		"fn":  "Push",
	})
	l.Debug("Pushing to nats")
	if d.Client == nil {
		l.Error("client is nil")
		return errors.New("client is nil")
	}
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		l.Errorf("error reading from reader: %v", err)
		return err
	}
	if err := d.Client.Publish(*d.Subject, bd); err != nil {
		l.Errorf("%+v", err)
		return err
	}
	l.Debug("Pushed to nats")
	return nil
}

func (d *NATS) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "nats",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up")
	if d.Client != nil {
		d.Client.Close()
	}
	l.Debug("Cleaned up")
	return nil
}

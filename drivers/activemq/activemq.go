package activemq

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"os"

	stomp "github.com/go-stomp/stomp/v3"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type ActiveMQ struct {
	Client  *stomp.Conn
	Address string
	Name    *string
	// TLS
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCert     *string
	TLSKey      *string
	TLSCA       *string
}

func (d *ActiveMQ) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "activemq",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	if os.Getenv(prefix+"ACTIVEMQ_ADDRESS") != "" {
		d.Address = os.Getenv(prefix + "ACTIVEMQ_ADDRESS")
	}
	if os.Getenv(prefix+"ACTIVEMQ_NAME") != "" {
		v := os.Getenv(prefix + "ACTIVEMQ_NAME")
		d.Name = &v
	}
	if os.Getenv(prefix+"ACTIVEMQ_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"ACTIVEMQ_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"ACTIVEMQ_TLS_INSECURE") != "" {
		v := os.Getenv(prefix+"ACTIVEMQ_TLS_INSECURE") == "true"
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"ACTIVEMQ_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "ACTIVEMQ_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"ACTIVEMQ_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "ACTIVEMQ_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	if os.Getenv(prefix+"ACTIVEMQ_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "ACTIVEMQ_TLS_CA_FILE")
		d.TLSCA = &v
	}
	l.Debug("Loaded environment")
	return nil
}

func (d *ActiveMQ) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "activemq",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.Address = *flags.ActiveMQAddress
	d.Name = flags.ActiveMQName
	d.EnableTLS = flags.ActiveMQEnableTLS
	d.TLSInsecure = flags.ActiveMQTLSInsecure
	d.TLSCert = flags.ActiveMQTLSCert
	d.TLSKey = flags.ActiveMQTLSKey
	d.TLSCA = flags.ActiveMQTLSCA
	l.Debug("Loaded flags")
	return nil
}

func (d *ActiveMQ) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "activemq",
		"fn":  "Init",
	})
	l.Debug("Initializing activemq driver")
	// create connection
	var err error
	var conn net.Conn
	if *d.EnableTLS {
		l.Debug("Creating TLS connection")
		tc, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
		if err != nil {
			return err
		}

		conn, err = tls.Dial("tcp", d.Address, tc)
		if err != nil {
			return err
		}
	} else {
		l.Debug("Creating non-TLS connection")
		conn, err = net.Dial("tcp", d.Address)
		if err != nil {
			return err
		}
	}
	d.Client, err = stomp.Connect(conn)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	l.Debug("Initialized activemq driver")
	return nil
}

func (d *ActiveMQ) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "activemq",
		"fn":  "Push",
	})
	l.Debug("Pushing message")
	jd, err := ioutil.ReadAll(r)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	err = d.Client.Send(*d.Name, "text/plain", jd, nil)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	l.Debug("Pushed message")
	return nil
}

func (d *ActiveMQ) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "activemq",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up")
	if err := d.Client.Disconnect(); err != nil {
		l.Errorf("%+v", err)
		return err
	}
	l.Debug("Cleaned up")
	return nil
}

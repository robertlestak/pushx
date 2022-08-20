package nsq

import (
	"io"
	"io/ioutil"
	"os"

	nsq "github.com/nsqio/go-nsq"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type NSQ struct {
	Client            *nsq.Producer
	NsqLookupdAddress *string
	NsqdAddress       *string
	Topic             *string
	// TLS
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCert     *string
	TLSKey      *string
	TLSCA       *string
}

func (d *NSQ) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "nsq",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	if os.Getenv(prefix+"NSQ_NSQLOOKUPD_ADDRESS") != "" {
		v := os.Getenv(prefix + "NSQ_NSQLOOKUPD_ADDRESS")
		d.NsqLookupdAddress = &v
	}
	if os.Getenv(prefix+"NSQ_NSQD_ADDRESS") != "" {
		v := os.Getenv(prefix + "NSQ_NSQD_ADDRESS")
		d.NsqdAddress = &v
	}
	if os.Getenv(prefix+"NSQ_TOPIC") != "" {
		v := os.Getenv(prefix + "NSQ_TOPIC")
		d.Topic = &v
	}
	if os.Getenv(prefix+"NSQ_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"NSQ_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"NSQ_TLS_INSECURE") != "" {
		v := os.Getenv(prefix+"NSQ_TLS_INSECURE") == "true"
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"NSQ_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "NSQ_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"NSQ_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "NSQ_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	if os.Getenv(prefix+"NSQ_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "NSQ_TLS_CA_FILE")
		d.TLSCA = &v
	}
	l.Debug("Loaded environment")
	return nil
}

func (d *NSQ) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "nsq",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.NsqLookupdAddress = flags.NSQNSQLookupdAddress
	d.NsqdAddress = flags.NSQNSQDAddress
	d.Topic = flags.NSQTopic
	d.EnableTLS = flags.NSQEnableTLS
	d.TLSInsecure = flags.NSQTLSSkipVerify
	d.TLSCert = flags.NSQCertFile
	d.TLSKey = flags.NSQKeyFile
	d.TLSCA = flags.NSQCAFile
	l.Debug("Loaded flags")
	return nil
}

func (d *NSQ) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "nsq",
		"fn":  "Init",
	})
	l.Debug("Initializing nsq driver")
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1
	if d.EnableTLS != nil && *d.EnableTLS {
		cfg.TlsV1 = true
		t, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
		if err != nil {
			l.Errorf("%+v", err)
			return err
		}
		cfg.TlsConfig = t
	}
	var addr string
	if d.NsqLookupdAddress != nil {
		addr = *d.NsqLookupdAddress
	} else {
		addr = *d.NsqdAddress
	}
	producer, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	switch os.Getenv("NSQ_LOG_LEVEL") {
	case "debug":
		producer.SetLoggerLevel(nsq.LogLevelDebug)
	case "info":
		producer.SetLoggerLevel(nsq.LogLevelInfo)
	case "warn":
		producer.SetLoggerLevel(nsq.LogLevelWarning)
	case "error":
		producer.SetLoggerLevel(nsq.LogLevelError)
	case "fatal":
		producer.SetLoggerLevel(nsq.LogLevelError)
	default:
		producer.SetLoggerLevel(nsq.LogLevelError)
	}
	d.Client = producer
	return nil
}

func (d *NSQ) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "nsq",
		"fn":  "Push",
	})
	l.Debug("Pushing message")
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	if err := d.Client.Publish(*d.Topic, bd); err != nil {
		l.Errorf("%+v", err)
		return err
	}
	return nil
}

func (d *NSQ) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "nsq",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up")
	d.Client.Stop()
	l.Debug("Cleaned up")
	return nil
}

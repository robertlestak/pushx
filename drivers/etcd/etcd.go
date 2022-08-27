package etcd

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	Client   *clientv3.Client
	Hosts    []string
	Username *string
	Password *string
	Key      string
	Limit    *int64
	// TLS
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCert     *string
	TLSKey      *string
	TLSCA       *string
}

func (d *Etcd) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "etcd",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"ETCD_HOSTS") != "" {
		d.Hosts = strings.Split(os.Getenv(prefix+"ETCD_HOSTS"), ",")
	}
	if os.Getenv(prefix+"ETCD_USERNAME") != "" {
		v := os.Getenv(prefix + "ETCD_USERNAME")
		d.Username = &v
	}
	if os.Getenv(prefix+"ETCD_PASSWORD") != "" {
		v := os.Getenv(prefix + "ETCD_PASSWORD")
		d.Password = &v
	}
	if os.Getenv(prefix+"ETCD_KEY") != "" {
		d.Key = os.Getenv(prefix + "ETCD_KEY")
	}
	if os.Getenv(prefix+"ETCD_TLS_ENABLE") != "" {
		v, err := strconv.ParseBool(os.Getenv(prefix + "ETCD_TLS_ENABLE"))
		if err != nil {
			l.Error(err)
			return err
		}
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"ETCD_TLS_INSECURE") != "" {
		v, err := strconv.ParseBool(os.Getenv(prefix + "ETCD_TLS_INSECURE"))
		if err != nil {
			l.Error(err)
			return err
		}
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"ETCD_TLS_CERT") != "" {
		v := os.Getenv(prefix + "ETCD_TLS_CERT")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"ETCD_TLS_KEY") != "" {
		v := os.Getenv(prefix + "ETCD_TLS_KEY")
		d.TLSKey = &v
	}
	if os.Getenv(prefix+"ETCD_TLS_CA") != "" {
		v := os.Getenv(prefix + "ETCD_TLS_CA")
		d.TLSCA = &v
	}
	return nil
}

func (d *Etcd) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "etcd",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.Username = flags.EtcdUsername
	d.Password = flags.EtcdPassword
	d.Hosts = strings.Split(*flags.EtcdHosts, ",")
	d.Key = *flags.EtcdKey
	d.EnableTLS = flags.EtcdTLSEnable
	d.TLSInsecure = flags.EtcdTLSInsecure
	d.TLSCert = flags.EtcdTLSCert
	d.TLSKey = flags.EtcdTLSKey
	d.TLSCA = flags.EtcdTLSCA
	return nil
}

func (d *Etcd) Init() error {
	l := log.WithFields(
		log.Fields{
			"pkg": "etcd",
			"fn":  "CreateFSSession",
		},
	)
	l.Debug("CreateFSSession")
	cfg := clientv3.Config{
		Endpoints:   d.Hosts,
		DialTimeout: 5 * time.Second,
	}
	if d.Username != nil && *d.Username != "" && d.Password != nil && *d.Password != "" {
		cfg.Username = *d.Username
		cfg.Password = *d.Password
	}
	if d.EnableTLS != nil && *d.EnableTLS {
		t, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
		if err != nil {
			l.Errorf("%+v", err)
			return err
		}
		cfg.TLS = t
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		l.Error(err)
		return err
	}
	d.Client = cli
	return nil
}

func (d *Etcd) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "etcd",
		"fn":  "Push",
	})
	l.Debug("Push")
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		l.Error(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var opts []clientv3.OpOption
	if d.Limit != nil && *d.Limit > 0 {
		opts = append(opts, clientv3.WithLimit(*d.Limit))
	}
	_, err = d.Client.Put(ctx, d.Key, string(bd), opts...)
	cancel()
	if err != nil {
		l.Error(err)
		return err
	}
	return nil
}

func (d *Etcd) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "etcd",
		"fn":  "Cleanup",
	})
	l.Debug("Cleanup")
	if err := d.Client.Close(); err != nil {
		l.Error(err)
		return err
	}
	return nil
}

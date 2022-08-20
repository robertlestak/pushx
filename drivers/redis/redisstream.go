package redis

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type RedisStream struct {
	Client    *redis.Client
	Host      string
	Port      string
	Password  string
	Key       string
	ValueKeys []string
	MessageID *string
	// TLS
	EnableTLS   *bool
	TLSInsecure *bool
	TLSCert     *string
	TLSKey      *string
	TLSCA       *string
}

func (d *RedisStream) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "redis",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"REDIS_HOST") != "" {
		d.Host = os.Getenv(prefix + "REDIS_HOST")
	}
	if os.Getenv(prefix+"REDIS_PORT") != "" {
		d.Port = os.Getenv(prefix + "REDIS_PORT")
	}
	if os.Getenv(prefix+"REDIS_PASSWORD") != "" {
		d.Password = os.Getenv(prefix + "REDIS_PASSWORD")
	}
	if os.Getenv(prefix+"REDIS_KEY") != "" {
		d.Key = os.Getenv(prefix + "REDIS_KEY")
	}
	if os.Getenv(prefix+"REDIS_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"REDIS_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"REDIS_TLS_INSECURE") != "" {
		v := os.Getenv(prefix+"REDIS_TLS_INSECURE") == "true"
		d.TLSInsecure = &v
	}
	if os.Getenv(prefix+"REDIS_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "REDIS_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"REDIS_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "REDIS_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	if os.Getenv(prefix+"REDIS_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "REDIS_TLS_CA_FILE")
		d.TLSCA = &v
	}
	if os.Getenv(prefix+"REDIS_MESSAGE_ID") != "" {
		v := os.Getenv(prefix + "REDIS_MESSAGE_ID")
		d.MessageID = &v
	}
	return nil
}

func (d *RedisStream) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "redis",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.Host = *flags.RedisHost
	d.Port = *flags.RedisPort
	d.Password = *flags.RedisPassword
	d.Key = *flags.RedisKey
	d.MessageID = flags.RedisMessageID
	d.EnableTLS = flags.RedisEnableTLS
	d.TLSInsecure = flags.RedisTLSSkipVerify
	d.TLSCert = flags.RedisCertFile
	d.TLSKey = flags.RedisKeyFile
	d.TLSCA = flags.RedisCAFile
	return nil
}

func (d *RedisStream) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "redis",
		"fn":  "Init",
	})
	l.Debug("Initializing redis stream driver")
	cfg := &redis.Options{
		Addr:        fmt.Sprintf("%s:%s", d.Host, d.Port),
		Password:    d.Password,
		DB:          0,
		DialTimeout: 30 * time.Second,
		ReadTimeout: 30 * time.Second,
	}
	if d.EnableTLS != nil && *d.EnableTLS {
		tc, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
		if err != nil {
			return err
		}
		cfg.TLSConfig = tc
	}
	d.Client = redis.NewClient(cfg)
	cmd := d.Client.Ping()
	if cmd.Err() != nil {
		l.Error("Failed to connect to redis")
		return cmd.Err()
	}
	l.Debug("Connected to redis")
	return nil
}

func (d *RedisStream) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "redis",
		"fn":  "Push",
	})
	l.Debug("Pushing to redis")
	var message map[string]interface{}
	if err := json.NewDecoder(r).Decode(&message); err != nil {
		l.Error("Failed to decode message")
		return err
	}
	if d.MessageID == nil || *d.MessageID == "" {
		v := "*"
		d.MessageID = &v
	}
	cmd := d.Client.XAdd(&redis.XAddArgs{
		Stream: d.Key,
		ID:     *d.MessageID,
		Values: message,
	})
	if cmd.Err() != nil {
		l.Error("Failed to push to redis")
		return cmd.Err()
	}
	l.Debug("Pushed to redis")
	return nil
}

func (d *RedisStream) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "redis",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up")
	if err := d.Client.Close(); err != nil {
		l.WithError(err).Error("Failed to close redis client")
		return err
	}
	return nil
}

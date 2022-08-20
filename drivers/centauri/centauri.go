package centauri

import (
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/robertlestak/centauri/pkg/agent"
	"github.com/robertlestak/centauri/pkg/keys"
	"github.com/robertlestak/centauri/pkg/message"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type Centauri struct {
	URL         string
	PublicKey   []byte
	MessageType string
	Filename    string
	Channel     *string
	Key         *string
}

func (d *Centauri) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "centauri",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	if os.Getenv(prefix+"CENTAURI_PEER_URL") != "" {
		d.URL = os.Getenv(prefix + "CENTAURI_PEER_URL")
	}
	if os.Getenv(prefix+"CENTAURI_CHANNEL") != "" {
		v := os.Getenv(prefix + "CENTAURI_CHANNEL")
		d.Channel = &v
	}
	if os.Getenv(prefix+"CENTAURI_PUBLIC_KEY") != "" {
		v := os.Getenv(prefix + "CENTAURI_PUBLIC_KEY")
		d.PublicKey = []byte(v)
	}
	if os.Getenv(prefix+"CENTAURI_MESSAGE_TYPE") != "" {
		d.MessageType = os.Getenv(prefix + "CENTAURI_MESSAGE_TYPE")
	}
	if os.Getenv(prefix+"CENTAURI_FILENAME") != "" {
		d.Filename = os.Getenv(prefix + "CENTAURI_FILENAME")
	}
	if os.Getenv(prefix+"CENTAURI_PUBLIC_KEY_BASE64") != "" {
		v := os.Getenv(prefix + "CENTAURI_PUBLIC_KEY_BASE64")
		dec, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			l.Errorf("error decoding base64: %v", err)
			return err
		}
		d.PublicKey = dec
	}
	return nil
}

func (d *Centauri) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "centauri",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.URL = *flags.CentauriPeerURL
	d.Channel = flags.CentauriChannel
	d.MessageType = *flags.CentauriMessageType
	d.Filename = *flags.CentauriFilename
	if flags.CentauriKeyBase64 != nil && *flags.CentauriKeyBase64 != "" {
		kd, err := base64.StdEncoding.DecodeString(*flags.CentauriKeyBase64)
		if err != nil {
			l.Errorf("error decoding base64: %v", err)
			return err
		}
		d.PublicKey = kd
	} else if flags.CentauriKey != nil && *flags.CentauriKey != "" {
		if flags.CentauriKey == nil || (flags.CentauriKey != nil && *flags.CentauriKey == "") {
			return errors.New("key required")
		}
		kd := []byte(*flags.CentauriKey)
		d.PublicKey = kd
		l.Debug("Loaded key", string(kd))
	}
	return nil
}

func (d *Centauri) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "centauri",
		"fn":  "Init",
	})
	l.Debug("Initializing centauri driver")
	if d.PublicKey == nil {
		l.Error("private key is nil")
		return errors.New("private key is nil")
	}
	agent.ServerAddrs = []string{d.URL}
	keys.AddKeyToPublicChain(d.PublicKey)
	return nil
}

func (d *Centauri) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "centauri",
		"fn":  "Push",
	})
	l.Debug("Pushing")
	if d.Channel == nil {
		defaultChannel := "default"
		d.Channel = &defaultChannel
	}
	if d.MessageType == "" {
		d.MessageType = "bytes"
	}
	m, err := message.CreateMessage(d.MessageType, d.Filename, *d.Channel, keys.PubKeyID(d.PublicKey), io.NopCloser(r))
	if err != nil {
		l.Errorf("error creating message: %v", err)
		return err
	}
	if err := agent.SendMessageThroughPeer(m); err != nil {
		l.Errorf("error sending message: %v", err)
		return err
	}
	l.Debug("Pushed")
	return nil
}

func (d *Centauri) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "centauri",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up")
	l.Debug("Cleaned up")
	return nil
}

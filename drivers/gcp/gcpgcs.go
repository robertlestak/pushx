package gcp

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type GCS struct {
	Client *storage.Client
	Bucket string
	Key    string
}

func (d *GCS) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"GCP_GCS_BUCKET") != "" {
		d.Bucket = os.Getenv(prefix + "GCP_GCS_BUCKET")
	}
	if os.Getenv(prefix+"GCP_GCS_KEY") != "" {
		d.Key = os.Getenv(prefix + "GCP_GCS_KEY")
	}
	return nil
}

func (d *GCS) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.Bucket = *flags.GCPGCSBucket
	d.Key = *flags.GCPGCSKey
	return nil
}

func (d *GCS) Init() error {
	l := log.WithFields(
		log.Fields{
			"pkg": "gcp",
			"fn":  "CreateGCPSession",
		},
	)
	l.Debug("CreateGCPSession")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	d.Client = client
	return err
}

func (d *GCS) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Push",
	})
	l.Debug("Push")
	if d.Bucket == "" {
		return fmt.Errorf("bucket is empty")
	}
	if d.Key == "" {
		return fmt.Errorf("key is empty")
	}
	ctx := context.Background()
	wc := d.Client.Bucket(d.Bucket).Object(d.Key).NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func (d *GCS) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Cleanup",
	})
	l.Debug("Cleanup")
	return nil
}

package gcp

import (
	"bytes"
	"context"
	"io"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type GCPPubSub struct {
	Client    *pubsub.Client
	ProjectID string
	TopicName string
}

func (d *GCPPubSub) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"GCP_PROJECT_ID") != "" {
		d.ProjectID = os.Getenv(prefix + "GCP_PROJECT_ID")
	}
	if os.Getenv(prefix+"GCP_TOPIC") != "" {
		d.TopicName = os.Getenv(prefix + "GCP_TOPIC")
	}
	return nil
}

func (d *GCPPubSub) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.ProjectID = *flags.GCPProjectID
	d.TopicName = *flags.GCPTopic
	return nil
}

func (d *GCPPubSub) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Init",
	})
	l.Debug("Initializing gcp pubsub driver")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, d.ProjectID)
	if err != nil {
		return err
	}
	d.Client = client
	return nil
}

func (d *GCPPubSub) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Push",
	})
	l.Debug("Pushing to gcp pubsub driver")
	ctx := context.Background()
	topic := d.Client.Topic(d.TopicName)
	defer topic.Stop()
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}
	msg := &pubsub.Message{
		Data: buf.Bytes(),
	}
	result := topic.Publish(ctx, msg)
	mid, err := result.Get(ctx)
	if err != nil {
		return err
	}
	l.WithField("mid", mid).Debug("Message published")
	return nil
}

func (d *GCPPubSub) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up gcp pubsub driver")
	return nil
}

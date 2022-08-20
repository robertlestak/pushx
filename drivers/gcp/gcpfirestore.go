package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type GCPFirestore struct {
	Client         *firestore.Client
	Collection     *string
	ID             *string
	FailCollection *string
	ProjectID      string
}

func (d *GCPFirestore) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"GCP_PROJECT_ID") != "" {
		d.ProjectID = os.Getenv(prefix + "GCP_PROJECT_ID")
	}
	if os.Getenv(prefix+"GCP_FIRESTORE_COLLECTION") != "" {
		v := os.Getenv(prefix + "GCP_FIRESTORE_COLLECTION")
		d.Collection = &v
	}
	if os.Getenv(prefix+"GCP_FIRESTORE_ID") != "" {
		v := os.Getenv(prefix + "GCP_FIRESTORE_ID")
		d.ID = &v
	}
	return nil
}

func (d *GCPFirestore) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.ProjectID = *flags.GCPProjectID
	d.Collection = flags.GCPFirestoreCollection
	d.ID = flags.GCPFirestoreID
	return nil
}

func (d *GCPFirestore) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Init",
	})
	l.Debug("Initializing gcp firestore driver")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, d.ProjectID)
	if err != nil {
		return err
	}
	d.Client = client
	return nil
}

func (d *GCPFirestore) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Push",
	})
	l.Debug("Creating document in gcp firestore driver")
	if *d.Collection == "" {
		return errors.New("no collection to create document in")
	}
	if *d.ID == "" {
		// generate new id
		v := uuid.New().String()
		d.ID = &v
	}
	ctx := context.Background()
	var message map[string]interface{}
	if err := json.NewDecoder(r).Decode(&message); err != nil {
		l.Error("Failed to decode message")
		return err
	}
	_, err := d.Client.Doc(*d.Collection+"/"+*d.ID).Create(ctx, message)
	if err != nil {
		l.Error(err)
		return err
	}
	return nil
}

func (d *GCPFirestore) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "gcp",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up gcp firestore driver")
	return nil
}

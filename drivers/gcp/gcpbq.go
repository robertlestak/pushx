package gcp

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/schema"

	log "github.com/sirupsen/logrus"
)

type BQ struct {
	Client    *bigquery.Client
	ProjectID string
	Query     *string
}

func (d *BQ) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "bq",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"GCP_BQ_QUERY") != "" {
		q := os.Getenv(prefix + "GCP_BQ_QUERY")
		d.Query = &q
	}
	if os.Getenv(prefix+"GCP_PROJECT_ID") != "" {
		d.ProjectID = os.Getenv(prefix + "GCP_PROJECT_ID")
	}
	return nil
}

func (d *BQ) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "bq",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	d.ProjectID = *flags.GCPProjectID
	if *flags.GCPBQQuery != "" {
		d.Query = flags.GCPBQQuery
	}
	return nil
}

func (d *BQ) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "bq",
		"fn":  "Init",
		"prj": d.ProjectID,
	})
	l.Debug("Initializing GCP_BQ client")
	var err error
	ctx := context.Background()
	c, err := bigquery.NewClient(ctx, d.ProjectID)
	if err != nil {
		return err
	}
	d.Client = c
	l.Debug("Connected to bq")
	return nil
}

func (d *BQ) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "bq",
		"fn":  "Push",
	})
	l.Debug("Handling failure in GCP_BQ")
	var err error
	if *d.Query == "" {
		return nil
	}
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if *d.Query != "" {
		q := schema.ReplaceParamsString(bd, *d.Query)
		d.Query = &q
		l.Debug("Query: " + *d.Query)
	}
	qry := d.Client.Query(*d.Query)
	_, err = qry.Read(context.Background())
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Pushed")
	return nil
}

func (d *BQ) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "bq",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up GCP_BQ")
	err := d.Client.Close()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Cleaned up")
	return nil
}

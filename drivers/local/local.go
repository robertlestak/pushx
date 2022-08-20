package local

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Local struct {
}

func (d *Local) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "local",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	return nil
}

func (d *Local) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "local",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	return nil
}

func (d *Local) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "local",
		"fn":  "Init",
	})
	l.Debug("Initializing local driver")
	return nil
}

func (d *Local) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "local",
		"fn":  "Push",
	})
	l.Debug("Pushing to local")
	// copy the input to stdout
	_, err := io.Copy(os.Stdout, r)
	if err != nil {
		l.WithError(err).Error("Copy")
		return err
	}
	return nil
}

func (d *Local) Cleanup() error {
	return nil
}

package fs

import (
	"io"
	"os"

	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type FS struct {
	Folder string
	Key    string
}

func (d *FS) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "fs",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"FS_KEY") != "" {
		d.Key = os.Getenv(prefix + "FS_KEY")
	}
	if os.Getenv(prefix+"FS_FOLDER") != "" {
		d.Folder = os.Getenv(prefix + "FS_FOLDER")
	}
	return nil
}

func (d *FS) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "fs",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.Folder = *flags.FSFolder
	d.Key = *flags.FSKey
	return nil
}

func (d *FS) Init() error {
	l := log.WithFields(
		log.Fields{
			"pkg": "fs",
			"fn":  "CreateFSSession",
		},
	)
	l.Debug("CreateFSSession")
	return nil
}

func (d *FS) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "fs",
		"fn":  "Push",
	})
	l.Debug("Push")
	// write the input to the file
	f, err := os.Create(d.Folder + "/" + d.Key)
	if err != nil {
		l.WithError(err).Error("Create")
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		l.WithError(err).Error("Copy")
		return err
	}
	return nil
}

func (d *FS) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "fs",
		"fn":  "Cleanup",
	})
	l.Debug("Cleanup")
	return nil
}

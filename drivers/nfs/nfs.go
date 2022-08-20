package nfs

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/go-nfs-client/nfs"
	"github.com/vmware/go-nfs-client/nfs/rpc"
)

type NFSMount struct {
	Mount  *nfs.Mount
	Target *nfs.Target
}

type NFS struct {
	Host   string
	Target string
	Folder string
	Key    string
	Client *NFSMount
}

func (d *NFS) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"NFS_HOST") != "" {
		d.Host = os.Getenv(prefix + "NFS_HOST")
	}
	if os.Getenv(prefix+"NFS_KEY") != "" {
		d.Key = os.Getenv(prefix + "NFS_KEY")
	}
	if os.Getenv(prefix+"NFS_FOLDER") != "" {
		d.Folder = os.Getenv(prefix + "NFS_FOLDER")
	}
	if os.Getenv(prefix+"NFS_TARGET") != "" {
		d.Target = os.Getenv(prefix + "NFS_TARGET")
	}
	return nil
}

func (d *NFS) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.Host = *flags.NFSHost
	d.Target = *flags.NFSTarget
	d.Folder = *flags.NFSFolder
	d.Key = *flags.NFSKey
	return nil
}

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.WithFields(log.Fields{
			"pkg": "nfs",
			"fn":  "hostname",
		}).Error("Failed to get hostname")
		return uuid.New().String()
	}
	return hostname
}

func (d *NFS) Init() error {
	l := log.WithFields(
		log.Fields{
			"pkg": "nfs",
			"fn":  "CreateNFSSession",
		},
	)
	l.Debug("CreateNFSSession")
	mount, err := nfs.DialMount(d.Host)
	if err != nil {
		log.Fatalf("unable to dial MOUNT service: %v", err)
	}
	auth := rpc.NewAuthUnix(hostname(), 1001, 1001)
	v, err := mount.Mount(d.Target, auth.Auth())
	if err != nil {
		log.Fatalf("unable to mount volume: %v", err)
	}
	d.Client = &NFSMount{
		Mount:  mount,
		Target: v,
	}
	return err
}

func (d *NFS) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "Push",
	})
	l.Debug("mvObject")
	if d.Key == "" {
		return errors.New("no key")
	}
	if d.Folder != "" {
		_, _, err := d.Client.Target.Lookup(d.Folder)
		if err != nil {
			// Create the bucket if it doesn't exist
			_, err = d.Client.Target.Mkdir(d.Folder, 0755)
			if err != nil {
				l.Errorf("mvObject error=%v", err)
				return err
			}
		}
	}
	w, err := d.Client.Target.OpenFile(path.Join(d.Folder, d.Key), 0644)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func (d *NFS) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "Cleanup",
	})
	l.Debug("Cleanup")
	if err := d.Client.Mount.Unmount(); err != nil {
		l.Errorf("%+v", err)
		return err
	}
	d.Client.Mount.Close()
	return nil
}

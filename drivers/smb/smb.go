package smb

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"

	"github.com/hirochachacha/go-smb2"
	"github.com/robertlestak/pushx/pkg/flags"
	log "github.com/sirupsen/logrus"
)

type SMBClient struct {
	Client *smb2.Session
	Conn   net.Conn
	Share  *smb2.Share
}

type SMB struct {
	Host     string
	Port     int
	Username *string
	Password *string
	Share    *string
	Key      string
	Client   *SMBClient
}

func (d *SMB) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "LoadEnv",
	})
	l.Debug("LoadEnv")
	if os.Getenv(prefix+"SMB_HOST") != "" {
		d.Host = os.Getenv(prefix + "SMB_HOST")
	}
	if os.Getenv(prefix+"SMB_PORT") != "" {
		pv, err := strconv.Atoi(os.Getenv(prefix + "SMB_PORT"))
		if err != nil {
			l.Errorf("%+v", err)
			return err
		}
		d.Port = pv
	}
	if os.Getenv(prefix+"SMB_USER") != "" {
		v := os.Getenv(prefix + "SMB_USER")
		d.Username = &v
	}
	if os.Getenv(prefix+"SMB_PASS") != "" {
		v := os.Getenv(prefix + "SMB_PASS")
		d.Password = &v
	}
	if os.Getenv(prefix+"SMB_SHARE") != "" {
		v := os.Getenv(prefix + "SMB_SHARE")
		d.Share = &v
	}
	if os.Getenv(prefix+"SMB_KEY") != "" {
		d.Key = os.Getenv(prefix + "SMB_KEY")
	}
	return nil
}

func (d *SMB) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "LoadFlags",
	})
	l.Debug("LoadFlags")
	d.Host = *flags.SMBHost
	d.Port = *flags.SMBPort
	d.Username = flags.SMBUser
	d.Password = flags.SMBPass
	d.Key = *flags.SMBKey
	d.Share = flags.SMBShare
	return nil
}

func (d *SMB) Init() error {
	l := log.WithFields(
		log.Fields{
			"pkg": "nfs",
			"fn":  "CreateNFSSession",
		},
	)
	l.Debug("CreateNFSSession")
	if d.Host == "" || d.Port == 0 || d.Username == nil || d.Password == nil || d.Share == nil {
		return errors.New("invalid SMB configuration")
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", d.Host, d.Port))
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	sd := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     *d.Username,
			Password: *d.Password,
		},
	}
	s, err := sd.Dial(conn)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	m, err := s.Mount(*d.Share)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	c := &SMBClient{
		Client: s,
		Conn:   conn,
		Share:  m,
	}
	d.Client = c
	return err
}

func (d *SMB) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "Push",
	})
	l.Debug("mvObject")
	if d.Key == "" {
		return errors.New("no key")
	}
	folder := path.Dir(d.Key)
	if folder != "" {
		in, err := d.Client.Share.Stat(folder)
		if err != nil || !in.IsDir() {
			err := d.Client.Share.MkdirAll(folder, 0755)
			if err != nil {
				l.Errorf("%+v", err)
				return err
			}
		}
	}
	f, err := d.Client.Share.Create(d.Key)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	return nil
}

func (d *SMB) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "nfs",
		"fn":  "Cleanup",
	})
	l.Debug("Cleanup")
	if d.Client == nil {
		return nil
	}
	d.Client.Share.Umount()
	d.Client.Client.Logoff()
	d.Client.Conn.Close()
	return nil
}

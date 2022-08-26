package scylla

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/schema"
	log "github.com/sirupsen/logrus"
)

type Scylla struct {
	Client      *gocql.Session
	Hosts       []string
	User        string
	Password    string
	LocalDC     *string
	Consistency string
	Keyspace    string
	Query       *schema.SqlQuery
}

func (d *Scylla) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "scylla",
		"fn":  "LoadEnv",
	})
	l.Debug("loading env")
	if os.Getenv(prefix+"SCYLLA_HOSTS") != "" {
		d.Hosts = strings.Split(os.Getenv(prefix+"SCYLLA_HOSTS"), ",")
	}
	if os.Getenv(prefix+"SCYLLA_KEYSPACE") != "" {
		d.Keyspace = os.Getenv(prefix + "SCYLLA_KEYSPACE")
	}
	if os.Getenv(prefix+"SCYLLA_USER") != "" {
		d.User = os.Getenv(prefix + "SCYLLA_USER")
	}
	if os.Getenv(prefix+"SCYLLA_PASSWORD") != "" {
		d.Password = os.Getenv(prefix + "SCYLLA_PASSWORD")
	}
	if os.Getenv(prefix+"SCYLLA_CONSISTENCY") != "" {
		d.Consistency = os.Getenv(prefix + "SCYLLA_CONSISTENCY")
	}
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if os.Getenv(prefix+"SCYLLA_QUERY") != "" {
		d.Query.Query = os.Getenv(prefix + "SCYLLA_QUERY")
	}
	if os.Getenv(prefix+"SCYLLA_LOCAL_DC") != "" {
		v := os.Getenv(prefix + "SCYLLA_LOCAL_DC")
		d.LocalDC = &v
	}
	if os.Getenv(prefix+"SCYLLA_PARAMS") != "" {
		for _, s := range strings.Split(os.Getenv(prefix+"SCYLLA_PARAMS"), ",") {
			d.Query.Params = append(d.Query.Params, s)
		}
	}
	return nil
}

func (d *Scylla) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "scylla",
		"fn":  "LoadFlags",
	})
	l.Debug("loading flags")
	var hosts []string
	if *flags.ScyllaHosts != "" {
		s := strings.Split(*flags.ScyllaHosts, ",")
		for _, v := range s {
			v = strings.TrimSpace(v)
			if v != "" {
				hosts = append(hosts, v)
			}
		}
	}
	var rps []any
	if *flags.ScyllaQueryParams != "" {
		s := strings.Split(*flags.ScyllaQueryParams, ",")
		for _, v := range s {
			rps = append(rps, v)
		}
	}
	d.Hosts = hosts
	d.User = *flags.ScyllaUser
	d.Password = *flags.ScyllaPassword
	d.Keyspace = *flags.ScyllaKeyspace
	d.Consistency = *flags.ScyllaConsistency
	d.LocalDC = flags.ScyllaLocalDC
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if *flags.ScyllaQuery != "" {
		d.Query.Query = *flags.ScyllaQuery
	}
	if len(rps) > 0 {
		d.Query.Params = rps
	}
	return nil
}

func (d *Scylla) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "scylla",
		"fn":  "Init",
	})
	l.Debug("Initializing scylla client")

	cluster := gocql.NewCluster(d.Hosts...)
	// parse consistency string
	consistencyLevel := gocql.ParseConsistency(d.Consistency)
	cluster.Consistency = consistencyLevel
	if d.Keyspace != "" {
		cluster.Keyspace = d.Keyspace
	}
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	if d.User != "" || d.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{Username: d.User, Password: d.Password}
	}
	fallback := gocql.RoundRobinHostPolicy()
	if d.LocalDC != nil && *d.LocalDC != "" {
		fallback = gocql.DCAwareRoundRobinPolicy(*d.LocalDC)
	}
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(fallback)
	if *d.LocalDC != "" {
		cluster.Consistency = gocql.LocalQuorum
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	d.Client = session
	return nil
}

func (d *Scylla) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "scylla",
		"fn":  "Push",
	})
	l.Debug("Pushing data")
	var err error
	if d.Query == nil || d.Query.Query == "" {
		return nil
	}
	bd, err := ioutil.ReadAll(r)
	if err != nil {
		l.Error(err)
		return err
	}
	d.Query.Params = schema.ReplaceParams(bd, d.Query.Params)
	l.Debugf("Executing query: %s %v", d.Query.Query, d.Query.Params)
	err = d.Client.Query(d.Query.Query, d.Query.Params...).Exec()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Query executed")
	return nil
}

func (d *Scylla) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "scylla",
		"fn":  "Cleanup",
	})
	l.Debug("cleaning up scylla")
	d.Client.Close()
	l.Debug("cleaned up scylla")
	return nil
}

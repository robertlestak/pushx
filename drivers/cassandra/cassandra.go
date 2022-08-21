package cassandra

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

type Cassandra struct {
	Client      *gocql.Session
	Hosts       []string
	User        string
	Password    string
	Consistency string
	Keyspace    string
	Query       *schema.SqlQuery
}

func (d *Cassandra) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "cassandra",
		"fn":  "LoadEnv",
	})
	l.Debug("loading env")
	if os.Getenv(prefix+"CASSANDRA_HOSTS") != "" {
		d.Hosts = strings.Split(os.Getenv(prefix+"CASSANDRA_HOSTS"), ",")
	}
	if os.Getenv(prefix+"CASSANDRA_KEYSPACE") != "" {
		d.Keyspace = os.Getenv(prefix + "CASSANDRA_KEYSPACE")
	}
	if os.Getenv(prefix+"CASSANDRA_USER") != "" {
		d.User = os.Getenv(prefix + "CASSANDRA_USER")
	}
	if os.Getenv(prefix+"CASSANDRA_PASSWORD") != "" {
		d.Password = os.Getenv(prefix + "CASSANDRA_PASSWORD")
	}
	if os.Getenv(prefix+"CASSANDRA_CONSISTENCY") != "" {
		d.Consistency = os.Getenv(prefix + "CASSANDRA_CONSISTENCY")
	}
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if os.Getenv(prefix+"CASSANDRA_QUERY") != "" {
		d.Query.Query = os.Getenv(prefix + "CASSANDRA_QUERY")
	}
	if os.Getenv(prefix+"CASSANDRA_PARAMS") != "" {
		for _, s := range strings.Split(os.Getenv(prefix+"CASSANDRA_PARAMS"), ",") {
			d.Query.Params = append(d.Query.Params, s)
		}
	}
	return nil
}

func (d *Cassandra) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "cassandra",
		"fn":  "LoadFlags",
	})
	l.Debug("loading flags")
	var hosts []string
	if *flags.CassandraHosts != "" {
		s := strings.Split(*flags.CassandraHosts, ",")
		for _, v := range s {
			v = strings.TrimSpace(v)
			if v != "" {
				hosts = append(hosts, v)
			}
		}
	}
	var rps []any
	if *flags.CassandraQueryParams != "" {
		s := strings.Split(*flags.CassandraQueryParams, ",")
		for _, v := range s {
			rps = append(rps, v)
		}
	}
	d.Hosts = hosts
	d.User = *flags.CassandraUser
	d.Password = *flags.CassandraPassword
	d.Keyspace = *flags.CassandraKeyspace
	d.Consistency = *flags.CassandraConsistency
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if *flags.CassandraQuery != "" {
		d.Query.Query = *flags.CassandraQuery
	}
	if len(rps) > 0 {
		d.Query.Params = rps
	}
	return nil
}

func (d *Cassandra) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "cassandra",
		"fn":  "Init",
	})
	l.Debug("Initializing cassandra client")

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
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	d.Client = session
	return nil
}

func (d *Cassandra) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "cassandra",
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

func (d *Cassandra) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "cassandra",
		"fn":  "Cleanup",
	})
	l.Debug("cleaning up cassandra")
	d.Client.Close()
	l.Debug("cleaned up cassandra")
	return nil
}

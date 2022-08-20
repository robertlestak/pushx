package cassandra

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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
	if os.Getenv(prefix+"CASSANDRA_QUERY") != "" {
		d.Query = &schema.SqlQuery{Query: os.Getenv(prefix + "CASSANDRA_QUERY")}
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
	d.User = *flags.CassandraUser
	d.Password = *flags.CassandraPassword
	d.Keyspace = *flags.CassandraKeyspace
	d.Consistency = *flags.CassandraConsistency
	if *flags.CassandraQuery != "" {
		rq := &schema.SqlQuery{
			Query:  *flags.CassandraQuery,
			Params: rps,
		}
		d.Query = rq
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

func extractMustacheKey(s string) string {
	l := log.WithFields(log.Fields{
		"pkg": "bq",
		"fn":  "extractMustacheKey",
	})
	l.Debug("Extracting mustache key")
	var key string
	for _, k := range strings.Split(s, "{{") {
		if strings.Contains(k, "}}") {
			key = strings.Split(k, "}}")[0]
			break
		}
	}
	return key
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
	for i, v := range d.Query.Params {
		sv := fmt.Sprintf("%s", v)
		if sv == "{{pushx_payload}}" {
			d.Query.Params[i] = bd
		} else if strings.Contains(sv, "{{") {
			key := extractMustacheKey(sv)
			d.Query.Params[i] = gjson.GetBytes(bd, key).String()
		}
	}
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

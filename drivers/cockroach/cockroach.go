package cockroach

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/schema"

	log "github.com/sirupsen/logrus"
)

type CockroachDB struct {
	Client      *sql.DB
	Host        string
	Port        int
	User        string
	Pass        string
	Db          string
	SslMode     string
	SSLRootCert *string
	SSLCert     *string
	SSLKey      *string
	RoutingID   *string
	Query       *schema.SqlQuery
}

func (d *CockroachDB) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "cockroach",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"COCKROACH_HOST") != "" {
		d.Host = os.Getenv(prefix + "COCKROACH_HOST")
	}
	if os.Getenv(prefix+"COCKROACH_PORT") != "" {
		pval, err := strconv.Atoi(os.Getenv(prefix + "COCKROACH_PORT"))
		if err != nil {
			return err
		}
		d.Port = pval
	}
	if os.Getenv(prefix+"COCKROACH_USER") != "" {
		d.User = os.Getenv(prefix + "COCKROACH_USER")
	}
	if os.Getenv(prefix+"COCKROACH_PASSWORD") != "" {
		d.Pass = os.Getenv(prefix + "COCKROACH_PASSWORD")
	}
	if os.Getenv(prefix+"COCKROACH_DATABASE") != "" {
		d.Db = os.Getenv(prefix + "COCKROACH_DATABASE")
	}
	if os.Getenv(prefix+"COCKROACH_SSL_MODE") != "" {
		d.SslMode = os.Getenv(prefix + "COCKROACH_SSL_MODE")
	}
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if os.Getenv(prefix+"COCKROACH_QUERY") != "" {
		d.Query.Query = os.Getenv(prefix + "COCKROACH_QUERY")
	}
	if os.Getenv(prefix+"COCKROACH_QUERY_PARAMS") != "" {
		p := strings.Split(os.Getenv(prefix+"COCKROACH_QUERY_PARAMS"), ",")
		for _, v := range p {
			d.Query.Params = append(d.Query.Params, v)
		}
	}
	if os.Getenv(prefix+"COCKROACH_ROUTING_ID") != "" {
		v := os.Getenv(prefix + "COCKROACH_ROUTING_ID")
		d.RoutingID = &v
	}
	if os.Getenv(prefix+"COCKROACH_TLS_ROOT_CERT") != "" {
		v := os.Getenv(prefix + "COCKROACH_TLS_ROOT_CERT")
		d.SSLRootCert = &v
	}
	if os.Getenv(prefix+"COCKROACH_TLS_CERT") != "" {
		v := os.Getenv(prefix + "COCKROACH_TLS_CERT")
		d.SSLCert = &v
	}
	if os.Getenv(prefix+"COCKROACH_TLS_KEY") != "" {
		v := os.Getenv(prefix + "COCKROACH_TLS_KEY")
		d.SSLKey = &v
	}
	return nil
}

func (d *CockroachDB) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "cockroach",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	pv, err := strconv.Atoi(*flags.CockroachDBPort)
	if err != nil {
		return err
	}
	var rps []any
	if *flags.CockroachDBQueryParams != "" {
		s := strings.Split(*flags.CockroachDBQueryParams, ",")
		for _, v := range s {
			rps = append(rps, v)
		}
	}
	d.Host = *flags.CockroachDBHost
	d.Port = pv
	d.User = *flags.CockroachDBUser
	d.Pass = *flags.CockroachDBPassword
	d.Db = *flags.CockroachDBDatabase
	d.SslMode = *flags.CockroachDBSSLMode
	d.RoutingID = flags.CockroachDBRoutingID
	d.SSLRootCert = flags.CockroachDBTLSRootCert
	d.SSLCert = flags.CockroachDBTLSCert
	d.SSLKey = flags.CockroachDBTLSKey
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if *flags.CockroachDBQuery != "" {
		d.Query.Query = *flags.CockroachDBQuery
	}
	if len(rps) > 0 {
		d.Query.Params = rps
	}
	return nil
}

func (d *CockroachDB) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "cockroach",
		"fn":  "Init",
	})
	l.Debug("Initializing cockroachdb client")
	var err error
	var connStr string = "postgresql://"
	var opts string
	if d.RoutingID != nil && *d.RoutingID != "" {
		opts = "&options=--cluster%3D" + *d.RoutingID
	}
	if d.User != "" && d.Pass != "" {
		connStr += fmt.Sprintf("%s:%s@%s:%d/%s",
			d.User, d.Pass, d.Host, d.Port, d.Db)
	} else if d.User != "" && d.Pass == "" {
		connStr += fmt.Sprintf("%s@%s:%d/%s",
			d.User, d.Host, d.Port, d.Db)
	}
	connStr += "?sslmode=" + d.SslMode
	if d.SSLRootCert != nil && *d.SSLRootCert != "" {
		connStr += "&sslrootcert=" + *d.SSLRootCert
	}
	if d.SSLCert != nil && *d.SSLCert != "" {
		connStr += "&sslcert=" + *d.SSLCert
	}
	if d.SSLKey != nil && *d.SSLKey != "" {
		connStr += "&sslkey=" + *d.SSLKey
	}
	connStr += opts
	l.Debugf("Connecting to %s", connStr)
	d.Client, err = sql.Open("postgres", connStr)
	if err != nil {
		l.Error(err)
		return err
	}
	err = d.Client.Ping()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Connected to cockroachdb")
	return nil
}

func (d *CockroachDB) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "cockroach",
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
	_, err = d.Client.Exec(d.Query.Query, d.Query.Params...)
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Pushed data")
	return nil
}

func (d *CockroachDB) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "cockroach",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up cockroachdb")
	err := d.Client.Close()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Cleaned up")
	return nil
}

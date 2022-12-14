package postgres

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

type Postgres struct {
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
	Query       *schema.SqlQuery
}

func (d *Postgres) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "postgres",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"PSQL_HOST") != "" {
		d.Host = os.Getenv(prefix + "PSQL_HOST")
	}
	if os.Getenv(prefix+"PSQL_PORT") != "" {
		pval, err := strconv.Atoi(os.Getenv(prefix + "PSQL_PORT"))
		if err != nil {
			return err
		}
		d.Port = pval
	}
	if os.Getenv(prefix+"PSQL_USER") != "" {
		d.User = os.Getenv(prefix + "PSQL_USER")
	}
	if os.Getenv(prefix+"PSQL_PASSWORD") != "" {
		d.Pass = os.Getenv(prefix + "PSQL_PASSWORD")
	}
	if os.Getenv(prefix+"PSQL_DATABASE") != "" {
		d.Db = os.Getenv(prefix + "PSQL_DATABASE")
	}
	if os.Getenv(prefix+"PSQL_SSL_MODE") != "" {
		d.SslMode = os.Getenv(prefix + "PSQL_SSL_MODE")
	}
	if os.Getenv(prefix+"PSQL_TLS_ROOT_CERT") != "" {
		v := os.Getenv(prefix + "PSQL_TLS_ROOT_CERT")
		d.SSLRootCert = &v
	}
	if os.Getenv(prefix+"PSQL_TLS_CERT") != "" {
		v := os.Getenv(prefix + "PSQL_TLS_CERT")
		d.SSLCert = &v
	}
	if os.Getenv(prefix+"PSQL_TLS_KEY") != "" {
		v := os.Getenv(prefix + "PSQL_TLS_KEY")
		d.SSLKey = &v
	}
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if os.Getenv(prefix+"PSQL_QUERY") != "" {
		d.Query.Query = os.Getenv(prefix + "PSQL_QUERY")
	}
	if os.Getenv(prefix+"PSQL_QUERY_PARAMS") != "" {
		p := strings.Split(os.Getenv(prefix+"PSQL_QUERY_PARAMS"), ",")
		for _, v := range p {
			d.Query.Params = append(d.Query.Params, v)
		}
	}
	return nil
}

func (d *Postgres) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "postgres",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	pv, err := strconv.Atoi(*flags.PsqlPort)
	if err != nil {
		return err
	}
	var rps []any
	if *flags.PsqlQueryParams != "" {
		s := strings.Split(*flags.PsqlQueryParams, ",")
		for _, v := range s {
			rps = append(rps, v)
		}
	}
	d.Host = *flags.PsqlHost
	d.Port = pv
	d.User = *flags.PsqlUser
	d.Pass = *flags.PsqlPassword
	d.Db = *flags.PsqlDatabase
	d.SslMode = *flags.PsqlSSLMode
	d.SSLRootCert = flags.PsqlTLSRootCert
	d.SSLCert = flags.PsqlTLSCert
	d.SSLKey = flags.PsqlTLSKey
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if *flags.PsqlQuery != "" {
		d.Query.Query = *flags.PsqlQuery
	}
	if len(rps) > 0 {
		d.Query.Params = rps
	}
	return nil
}

func (d *Postgres) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "postgres",
		"fn":  "Init",
	})
	l.Debug("Initializing psql client")
	var err error
	var opts string
	var connStr string = "postgresql://"
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
	l.Debug("Connected to psql")
	return nil
}

func (d *Postgres) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "postgres",
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

func (d *Postgres) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "postgres",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up psql")
	err := d.Client.Close()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Cleaned up")
	return nil
}

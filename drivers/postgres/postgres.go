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
	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

type Postgres struct {
	Client  *sql.DB
	Host    string
	Port    int
	User    string
	Pass    string
	Db      string
	SslMode string
	Query   *schema.SqlQuery
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
	if os.Getenv(prefix+"PSQL_QUERY") != "" {
		d.Query = &schema.SqlQuery{
			Query: os.Getenv(prefix + "PSQL_QUERY"),
		}
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
	if *flags.PsqlQuery != "" {
		rq := &schema.SqlQuery{
			Query:  *flags.PsqlQuery,
			Params: rps,
		}
		d.Query = rq
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
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Pass, d.Host, d.Port, d.Db, d.SslMode)
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
	// loop through params and if we find {{key}}, replace it with the key
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

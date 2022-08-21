package mssql

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/schema"
	log "github.com/sirupsen/logrus"
)

type MSSql struct {
	Client *sql.DB
	Host   string
	Port   int
	User   string
	Pass   string
	Db     string
	Query  *schema.SqlQuery
}

func (d *MSSql) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "mssql",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"MSSQL_HOST") != "" {
		d.Host = os.Getenv(prefix + "MSSQL_HOST")
	}
	if os.Getenv(prefix+"MSSQL_PORT") != "" {
		pv, err := strconv.Atoi(os.Getenv(prefix + "MSSQL_PORT"))
		if err != nil {
			return err
		}
		d.Port = pv
	}
	if os.Getenv(prefix+"MSSQL_USER") != "" {
		d.User = os.Getenv(prefix + "MSSQL_USER")
	}
	if os.Getenv(prefix+"MSSQL_PASSWORD") != "" {
		d.Pass = os.Getenv(prefix + "MSSQL_PASSWORD")
	}
	if os.Getenv(prefix+"MSSQL_DATABASE") != "" {
		d.Db = os.Getenv(prefix + "MSSQL_DATABASE")
	}
	if os.Getenv(prefix+"MSSQL_QUERY") != "" {
		d.Query = &schema.SqlQuery{
			Query: os.Getenv(prefix + "MSSQL_QUERY"),
		}
	}
	if os.Getenv(prefix+"MSSQL_QUERY_PARAMS") != "" {
		p := strings.Split(os.Getenv(prefix+"MSSQL_QUERY_PARAMS"), ",")
		for _, v := range p {
			d.Query.Params = append(d.Query.Params, v)
		}
	}
	return nil
}

func (d *MSSql) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "mssql",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	pv, err := strconv.Atoi(*flags.MSSqlPort)
	if err != nil {
		return err
	}
	var rps []any
	if *flags.MSSqlQueryParams != "" {
		s := strings.Split(*flags.MSSqlQueryParams, ",")
		for _, v := range s {
			rps = append(rps, v)
		}
	}
	d.Host = *flags.MSSqlHost
	d.Port = pv
	d.User = *flags.MSSqlUser
	d.Pass = *flags.MSSqlPassword
	d.Db = *flags.MSSqlDatabase
	if *flags.MSSqlQuery != "" {
		rq := &schema.SqlQuery{
			Query:  *flags.MSSqlQuery,
			Params: rps,
		}
		d.Query = rq
	}
	return nil
}

func (d *MSSql) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "mssql",
		"fn":  "Init",
	})
	l.Debug("Initializing mssql client")
	var err error
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", d.Host, d.User, d.Pass, d.Port, d.Db)
	l.Debug("Connecting to mssql: ", connStr)
	d.Client, err = sql.Open("mssql", connStr)
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Initialized mssql client")
	// ping the database to check if it is alive
	err = d.Client.Ping()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Pinged mssql client")
	return nil
}

func (d *MSSql) Push(r io.Reader) error {
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

func (d *MSSql) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "mssql",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up mssql client")
	err := d.Client.Close()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Cleaned up mssql client")
	return nil
}

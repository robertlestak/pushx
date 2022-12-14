package mysql

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/schema"
	log "github.com/sirupsen/logrus"
)

type Mysql struct {
	Client *sql.DB
	Host   string
	Port   int
	User   string
	Pass   string
	Db     string
	Query  *schema.SqlQuery
}

func (d *Mysql) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "mysql",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment variables")
	if os.Getenv(prefix+"MYSQL_HOST") != "" {
		d.Host = os.Getenv(prefix + "MYSQL_HOST")
	}
	if os.Getenv(prefix+"MYSQL_PORT") != "" {
		pv, err := strconv.Atoi(os.Getenv(prefix + "MYSQL_PORT"))
		if err != nil {
			return err
		}
		d.Port = pv
	}
	if os.Getenv(prefix+"MYSQL_USER") != "" {
		d.User = os.Getenv(prefix + "MYSQL_USER")
	}
	if os.Getenv(prefix+"MYSQL_PASSWORD") != "" {
		d.Pass = os.Getenv(prefix + "MYSQL_PASSWORD")
	}
	if os.Getenv(prefix+"MYSQL_DATABASE") != "" {
		d.Db = os.Getenv(prefix + "MYSQL_DATABASE")
	}
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if os.Getenv(prefix+"MYSQL_QUERY") != "" {
		d.Query.Query = os.Getenv(prefix + "MYSQL_QUERY")
	}
	if os.Getenv(prefix+"MYSQL_QUERY_PARAMS") != "" {
		p := strings.Split(os.Getenv(prefix+"MYSQL_QUERY_PARAMS"), ",")
		for _, v := range p {
			d.Query.Params = append(d.Query.Params, v)
		}
	}
	return nil
}

func (d *Mysql) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "mysql",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	pv, err := strconv.Atoi(*flags.MysqlPort)
	if err != nil {
		return err
	}
	var rps []any
	if *flags.MysqlQueryParams != "" {
		s := strings.Split(*flags.MysqlQueryParams, ",")
		for _, v := range s {
			rps = append(rps, v)
		}
	}
	d.Host = *flags.MysqlHost
	d.Port = pv
	d.User = *flags.MysqlUser
	d.Pass = *flags.MysqlPassword
	d.Db = *flags.MysqlDatabase
	if d.Query == nil {
		d.Query = &schema.SqlQuery{}
	}
	if *flags.MysqlQuery != "" {
		d.Query.Query = *flags.MysqlQuery
	}
	if len(rps) > 0 {
		d.Query.Params = rps
	}
	return nil
}

func (d *Mysql) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "mysql",
		"fn":  "Init",
	})
	l.Debug("Initializing mysql client")
	var err error
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Pass, d.Host, d.Port, d.Db)
	l.Debug("Connecting to mysql: ", connStr)
	d.Client, err = sql.Open("mysql", connStr)
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Initialized mysql client")
	// ping the database to check if it is alive
	err = d.Client.Ping()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Pinged mysql client")
	return nil
}

func (d *Mysql) Push(r io.Reader) error {
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

func (d *Mysql) Cleanup() error {
	l := log.WithFields(log.Fields{
		"pkg": "mysql",
		"fn":  "Cleanup",
	})
	l.Debug("Cleaning up mysql client")
	err := d.Client.Close()
	if err != nil {
		l.Error(err)
		return err
	}
	l.Debug("Cleaned up mysql client")
	return nil
}

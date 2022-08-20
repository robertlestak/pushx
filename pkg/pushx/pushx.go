package pushx

import (
	"io"
	"os"
	"strings"

	"github.com/robertlestak/pushx/pkg/drivers"
	log "github.com/sirupsen/logrus"
)

type PushX struct {
	DriverName drivers.DriverName `json:"driverName"`
	Driver     drivers.Driver     `json:"driver"`
	InputStr   string             `json:"in"`
	InputFile  string             `json:"inFile"`
	Input      io.Reader          `json:"-"`
	OutputFile string             `json:"outFile"`
	Output     io.Writer          `json:"-"`
}

func (j *PushX) Init(envKeyPrefix string) error {
	l := log.WithFields(log.Fields{
		"fn": "Init",
	})
	l.Debug("Init")
	if j.DriverName == "" {
		l.Error("no driver specified")
		return drivers.ErrDriverNotFound
	}
	l.Debug("driver specified")
	j.Driver = drivers.GetDriver(j.DriverName)
	if j.Driver == nil {
		l.Error("driver not found")
		return drivers.ErrDriverNotFound
	}
	if err := j.Driver.LoadFlags(); err != nil {
		l.WithError(err).Error("LoadFlags")
		return err
	}
	if err := j.Driver.LoadEnv(envKeyPrefix); err != nil {
		l.WithError(err).Error("LoadEnv")
		return err
	}
	if err := j.Driver.Init(); err != nil {
		l.WithError(err).Error("Init")
		return err
	}
	l.Debug("driver initialized")
	if j.InputStr != "" {
		l.Debug("input string specified")
		j.Input = strings.NewReader(j.InputStr)
	} else if j.InputFile == "-" {
		l.Debug("input is stdin")
		j.Input = os.Stdin
	} else {
		l.Debug("input is file")
		var err error
		j.Input, err = os.Open(j.InputFile)
		if err != nil {
			l.WithError(err).Error("Open")
			return err
		}
	}
	return nil
}

func (j *PushX) output() error {
	l := log.WithFields(log.Fields{
		"fn":     "output",
		"driver": j.DriverName,
	})
	l.Debug("output")
	if j.OutputFile == "-" {
		l.Debug("output is stdout")
		j.Output = os.Stdout
	} else if j.OutputFile != "" {
		l.Debug("output is file")
		var err error
		j.Output, err = os.Create(j.OutputFile)
		if err != nil {
			log.WithError(err).Error("Create")
			return err
		}
	} else {
		l.Debug("no output")
	}
	return nil
}

func (j *PushX) Push() error {
	l := log.WithFields(log.Fields{
		"fn":     "Push",
		"driver": j.DriverName,
	})
	var in io.Reader
	if j.OutputFile == "" {
		l.Debug("no output")
		in = j.Input
	} else {
		if err := j.output(); err != nil {
			l.WithError(err).Error("output")
			return err
		}
		in = io.TeeReader(j.Input, j.Output)
	}
	err := j.Driver.Push(in)
	if err != nil {
		l.Error("push error:", err)
		return err
	}
	l.Debug("work pushed")
	return nil
}

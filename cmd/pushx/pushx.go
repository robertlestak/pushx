package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/robertlestak/pushx/pkg/drivers"
	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/pushx"
	log "github.com/sirupsen/logrus"
)

var (
	Version      = "dev"
	AppName      = "pushx"
	EnvKeyPrefix = fmt.Sprintf("%s_", strings.ToUpper(AppName))
)

func init() {
	ll, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}

func printVersion() {
	fmt.Printf(AppName+" version %s\n", Version)
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", AppName)
	flags.FlagSet.PrintDefaults()
}

func LoadEnv(prefix string) error {
	if os.Getenv(prefix+"DRIVER") != "" {
		d := os.Getenv(prefix + "DRIVER")
		flags.Driver = &d
	}
	if os.Getenv(prefix+"INPUT_STR") != "" {
		i := os.Getenv(prefix + "INPUT_STR")
		flags.InputStr = &i
	}
	if os.Getenv(prefix+"INPUT_FILE") != "" {
		i := os.Getenv(prefix + "INPUT_FILE")
		flags.InputFile = &i
	}
	return nil
}

func run(j *pushx.PushX) error {
	l := log.WithFields(log.Fields{
		"app": AppName,
		"fn":  "run",
	})
	l.Debug("start")
	if err := j.Push(); err != nil {
		l.Errorf("failed to do work: %s", err)
		return err
	}
	l.Debug("done")
	return nil
}

func cleanup(j *pushx.PushX) error {
	l := log.WithFields(log.Fields{
		"app": AppName,
		"fn":  "cleanup",
	})
	l.Debug("cleanup")
	if err := j.Driver.Cleanup(); err != nil {
		l.Error(err)
		return err
	}
	return nil
}

func main() {
	l := log.WithFields(log.Fields{
		"app": AppName,
		"fn":  "main",
	})
	l.Debug("start")
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			printVersion()
			os.Exit(0)
		}
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			printUsage()
			os.Exit(0)
		}
	}
	flags.FlagSet.Parse(os.Args[1:])
	if err := LoadEnv(EnvKeyPrefix); err != nil {
		l.Error(err)
		os.Exit(1)
	}
	l.Debug("parsed flags")
	j := &pushx.PushX{
		DriverName: drivers.DriverName(*flags.Driver),
		InputStr:   *flags.InputStr,
		InputFile:  *flags.InputFile,
	}
	if err := j.Init(EnvKeyPrefix); err != nil {
		l.WithError(err).Error("InitDriver")
		os.Exit(1)
	}
	l.Debug("initialized driver")
	if err := run(j); err != nil {
		l.WithError(err).Error("run")
		os.Exit(1)
	}
	if err := cleanup(j); err != nil {
		l.WithError(err).Error("cleanup")
		os.Exit(1)
	}
	l.Debug("exited")
}

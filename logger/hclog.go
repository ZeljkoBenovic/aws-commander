package logger

import (
	"github.com/Trapesys/aws-commander/conf"
	"github.com/hashicorp/go-hclog"
	"os"
)

type hlog struct {
	log hclog.Logger
}

func (h *hlog) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		h.log.Error(msg)
		return
	}

	h.log.Error(msg, args...)
}

func (h *hlog) Warn(msg string, args ...interface{}) {
	if len(args) == 0 {
		h.log.Warn(msg)
		return
	}

	h.log.Warn(msg, args...)
}

func (h *hlog) Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		h.log.Info(msg)
		return
	}

	h.log.Info(msg, args...)
}

func (h *hlog) Debug(msg string, args ...interface{}) {
	if len(args) == 0 {
		h.log.Debug(msg)
		return
	}

	h.log.Debug(msg, args...)
}

func (h *hlog) Fatalln(msg string, args ...interface{}) {
	if len(args) == 0 {
		h.log.Error(msg)
	} else {
		h.log.Error(msg, args...)
	}

	os.Exit(1)
}

func (h *hlog) Named(name string) Logger {
	h.log = h.log.Named(name)
	return h
}

func New(conf conf.Config) Logger {
	return &hlog{
		log: hclog.New(&hclog.LoggerOptions{
			Name:  "aws-commander",
			Level: hclog.LevelFromString(conf.LogLevel),
			Color: hclog.AutoColor,
		}),
	}
}

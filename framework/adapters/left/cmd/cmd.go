package cmd

import (
	"flag"
	"github.com/Trapesys/aws-commander/framework/adapters/types/cmd"
	"github.com/Trapesys/aws-commander/framework/ports"
	"github.com/hashicorp/go-hclog"
	"os"
	"strings"
)

type Adapter struct {
	flags  *cmd.Flags
	logger hclog.Logger

	buffInstanceFlag string
}

func NewAdapter() ports.ICmd {
	return &Adapter{
		flags: &cmd.Flags{
			AwsZone:            new(string),
			BashScriptLocation: new(string),
			LogLevel:           new(string),
			OutputLocation:     new(string),
			FreeFormCmd:        new(string),
			AwsProfile:         new(string),
		},
	}
}

func (a *Adapter) GetFlags() cmd.Flags {
	flag.StringVar(a.flags.AwsZone, "aws-zone", "eu-central-1", "aws zone where instances reside")
	flag.StringVar(a.flags.BashScriptLocation, "script", "", "the location of the script to run")
	flag.StringVar(&a.buffInstanceFlag, "instances", "", "instance IDs, separated by comma (,)")
	flag.StringVar(a.flags.LogLevel, "log-level", "info", "log output level")
	flag.StringVar(a.flags.OutputLocation, "output", "", "the location of file to write json output "+
		"(default: output to console)")
	flag.StringVar(a.flags.FreeFormCmd, "cmd", "", "freeform command, a single line bash command to be executed")
	flag.StringVar(a.flags.AwsProfile, "aws-profile", "default", "aws credentials profile")
	flag.Parse()

	if a.buffInstanceFlag == "" {
		a.logger.Error("instance IDs not defined")
		flag.PrintDefaults()

		os.Exit(1)
	}

	a.flags.InstanceIDs = append(a.flags.InstanceIDs, strings.Split(a.buffInstanceFlag, ",")...)

	return *a.flags
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.ICmd {
	a.logger = logger.Named("cmd")

	return a
}

package app

import (
	"fmt"
	"github.com/ZeljkoBenovic/aws-commander/framework/adapters/types/cmd"
	fports "github.com/ZeljkoBenovic/aws-commander/framework/ports"
	"github.com/ZeljkoBenovic/aws-commander/internal/ports"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
	"os"
)

type Adapter struct {
	coreAdapter    ports.CorePort
	localFSAdapter fports.ILocalFS
	ssmAdapter     fports.ISSMPort
	cmdAdapter     fports.ICmd
	baseLogger     hclog.Logger

	flags cmd.Flags
}

func (a Adapter) getAWSOptions() session.Options {
	// use env vars if set
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" &&
		os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		a.baseLogger.Info("using aws secrets found in env vars")

		return session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("AWS_ACCESS_KEY_ID"),
					os.Getenv("AWS_SECRET_ACCESS_KEY"),
					os.Getenv("AWS_SESSION_TOKEN"),
				),
				Region: a.flags.AwsZone,
			},
		}
	}

	a.baseLogger.Info("reading aws credentials from", "profile", *a.flags.AwsProfile)
	// otherwise, use profile from aws credential file
	return session.Options{
		Config: aws.Config{
			Region: a.flags.AwsZone,
		},
		Profile: *a.flags.AwsProfile,
	}
}

func (a *Adapter) Init() ports.IApp {
	/// CMD ///
	a.flags = a.cmdAdapter.WithLogger(a.baseLogger).GetFlags()
	//////////
	// set log level
	a.baseLogger.SetLevel(hclog.LevelFromString(*a.flags.LogLevel))

	// CORE //
	// init core
	a.coreAdapter.WithAWSSessionOptions(a.getAWSOptions()).WithLogger(a.baseLogger)
	/////////

	// SSM //
	// init ssm session
	a.ssmAdapter.WithAWSSession(a.coreAdapter.GetAWSSession()).WithLogger(a.baseLogger)

	// init ssm instances
	a.ssmAdapter.WithInstances(a.flags.InstanceIDs)
	//init ssm commands
	if *a.flags.BashScriptLocation != "" {
		a.ssmAdapter.WithCommands(
			a.localFSAdapter.
				WithLogger(a.baseLogger).
				ReadBashScript(*a.flags.BashScriptLocation),
		)
	} else if *a.flags.FreeFormCmd != "" {
		a.baseLogger.Debug("freeform command found", "cmd", *a.flags.FreeFormCmd)
		a.ssmAdapter.WithFreeFormCommand(*a.flags.FreeFormCmd)
	} else {
		a.baseLogger.Error("could not find a command to execute, " +
			"script location or cmd flag must be defined")
		os.Exit(1)
	}
	/////////////

	return a
}

func NewAdapter(
	core ports.CorePort,
	local fports.ILocalFS,
	ssm fports.ISSMPort,
	cmd fports.ICmd) ports.IApp {
	return &Adapter{
		coreAdapter:    core,
		localFSAdapter: local,
		ssmAdapter:     ssm,
		cmdAdapter:     cmd,
	}
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.IApp {
	a.baseLogger = logger

	return a
}

func (a Adapter) RunCommand() error {
	cmdResult := a.ssmAdapter.RunCommand()
	if *a.flags.OutputLocation != "" {
		if err := a.localFSAdapter.WriteRunCommandOutput(cmdResult, *a.flags.OutputLocation); err != nil {
			return fmt.Errorf("could not write command result to file: %w", err)
		}

		a.baseLogger.Info("command output written to file", "filename", *a.flags.OutputLocation)
		os.Exit(0)
	}

	// output to console
	fmt.Printf("%+v\n", cmdResult)

	return nil
}

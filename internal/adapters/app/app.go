package app

import (
	"fmt"
	"os"

	"github.com/Trapesys/aws-commander/framework/adapters/types/cmd"
	fports "github.com/Trapesys/aws-commander/framework/ports"
	"github.com/Trapesys/aws-commander/internal/ports"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
)

type Adapter struct {
	coreAdapter    ports.CorePort
	localFSAdapter fports.ILocalFS
	ssmAdapter     fports.ISSMPort
	cmdAdapter     fports.ICmd
	baseLogger     hclog.Logger

	flags cmd.Flags
}

func (a *Adapter) getAWSOptions() session.Options {
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
				Region: a.flags.AwsZone.ValueString,
			},
		}
	}

	// otherwise, use profile from aws credential file
	awsProfile := os.Getenv("AWS_PROFILE")
	if awsProfile == "" {
		awsProfile = *a.flags.AwsProfile.ValueString
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = *a.flags.AwsZone.ValueString
	}

	a.baseLogger.Info("reading aws credentials from", "profile", awsProfile, "region", awsRegion)

	return session.Options{
		Config: aws.Config{
			Region: &awsRegion,
		},
		Profile:           awsProfile,
		SharedConfigState: session.SharedConfigEnable,
	}
}

func (a *Adapter) Init() ports.IApp {
	/// CMD ///
	a.flags = a.cmdAdapter.WithLogger(a.baseLogger).GetFlags()
	//////////
	// set log level
	a.baseLogger.SetLevel(hclog.LevelFromString(*a.flags.LogLevel.ValueString))

	// CORE //
	// init core
	a.coreAdapter.WithAWSSessionOptions(a.getAWSOptions()).WithLogger(a.baseLogger)
	/////////

	// SSM //
	// init ssm session
	a.ssmAdapter.WithAWSSession(a.coreAdapter.GetAWSSession()).WithLogger(a.baseLogger)

	// MODE //
	// init command mode ( bash or ansible )
	a.ssmAdapter.WithMode(*a.flags.Mode.ValueString)

	// init ssm instances
	a.ssmAdapter.WithInstances(a.flags.InstanceIDs.ValueStringArr)

	// init ssm commands
	switch *a.flags.Mode.ValueString {
	case "bash":
		a.baseLogger.Debug("mode set to Bash")

		if *a.flags.BashScriptLocation.ValueString != "" {
			a.baseLogger.Debug("running Bash script", "script_location", *a.flags.BashScriptLocation.ValueString)
			a.ssmAdapter.WithCommands(
				a.localFSAdapter.
					WithLogger(a.baseLogger).
					ReadBashScript(*a.flags.BashScriptLocation.ValueString),
			)

			break
		}

		if *a.flags.FreeFormCmd.ValueString != "" {
			a.baseLogger.Debug("freeform command found", "cmd", *a.flags.FreeFormCmd.ValueString)
			a.ssmAdapter.WithFreeFormCommand(*a.flags.FreeFormCmd.ValueString)

			break
		}

		a.baseLogger.Error("when Bash mode is selected, script location and/or freeform command can't be empty")

	case "ansible":
		a.baseLogger.Debug("mode set to Ansible")

		// TODO: implement PlaybookURL and ExtraVars
		isDryRun := "False"
		if *a.flags.AnsibleDryRun.ValueBool {
			isDryRun = "True"
		}

		a.ssmAdapter.WithAnsiblePlaybook(&fports.AnsiblePlaybookOpts{
			Playbook:    a.localFSAdapter.ReadAnsiblePlaybook(*a.flags.AnsiblePlaybook.ValueString),
			PlaybookURL: "",
			ExtraVars:   "",
			Check:       isDryRun,
		})
	default:
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

func (a *Adapter) RunCommand() error {
	cmdResult := a.ssmAdapter.RunCommand()
	if *a.flags.OutputLocation.ValueString != "" {
		if err := a.localFSAdapter.WriteRunCommandOutput(cmdResult, *a.flags.OutputLocation.ValueString); err != nil {
			return fmt.Errorf("could not write command result to file: %w", err)
		}

		a.baseLogger.Info("command output written to file", "filename", *a.flags.OutputLocation.ValueString)
		os.Exit(0)
	}

	// output to console
	fmt.Printf("%+v\n", cmdResult)

	return nil
}

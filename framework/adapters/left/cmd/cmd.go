package cmd

import (
	"flag"
	"os"
	"strings"

	"github.com/Trapesys/aws-commander/framework/adapters/types/cmd"
	"github.com/Trapesys/aws-commander/framework/ports"
	"github.com/hashicorp/go-hclog"
)

type Adapter struct {
	logger hclog.Logger

	buffInstanceFlag string
}

func NewAdapter() ports.ICmd {
	return &Adapter{}
}

func (a *Adapter) GetFlags() cmd.Flags {
	flag.StringVar(cmd.UserFlags.AwsZone.ValueString,
		cmd.UserFlags.AwsZone.Name,
		cmd.UserFlags.AwsZone.DefaultString,
		cmd.UserFlags.AwsZone.Usage,
	)
	flag.StringVar(cmd.UserFlags.BashScriptLocation.ValueString,
		cmd.UserFlags.BashScriptLocation.Name,
		cmd.UserFlags.BashScriptLocation.DefaultString,
		cmd.UserFlags.BashScriptLocation.Usage,
	)
	flag.StringVar(&a.buffInstanceFlag,
		cmd.UserFlags.InstanceIDs.Name,
		cmd.UserFlags.InstanceIDs.DefaultString,
		cmd.UserFlags.InstanceIDs.Usage,
	)
	flag.StringVar(cmd.UserFlags.LogLevel.ValueString,
		cmd.UserFlags.LogLevel.Name,
		cmd.UserFlags.LogLevel.DefaultString,
		cmd.UserFlags.LogLevel.Usage,
	)
	flag.StringVar(cmd.UserFlags.OutputLocation.ValueString,
		cmd.UserFlags.OutputLocation.Name,
		cmd.UserFlags.OutputLocation.DefaultString,
		cmd.UserFlags.OutputLocation.Usage,
	)
	flag.StringVar(cmd.UserFlags.FreeFormCmd.ValueString,
		cmd.UserFlags.FreeFormCmd.Name,
		cmd.UserFlags.FreeFormCmd.DefaultString,
		cmd.UserFlags.FreeFormCmd.Usage,
	)
	flag.StringVar(cmd.UserFlags.AwsProfile.ValueString,
		cmd.UserFlags.AwsProfile.Name,
		cmd.UserFlags.AwsProfile.DefaultString,
		cmd.UserFlags.AwsProfile.Usage,
	)
	flag.StringVar(cmd.UserFlags.Mode.ValueString,
		cmd.UserFlags.Mode.Name,
		cmd.UserFlags.Mode.DefaultString,
		cmd.UserFlags.Mode.Usage,
	)
	flag.StringVar(cmd.UserFlags.AnsiblePlaybook.ValueString,
		cmd.UserFlags.AnsiblePlaybook.Name,
		cmd.UserFlags.AnsiblePlaybook.DefaultString,
		cmd.UserFlags.AnsiblePlaybook.Usage,
	)
	flag.BoolVar(cmd.UserFlags.AnsibleDryRun.ValueBool,
		cmd.UserFlags.AnsibleDryRun.Name,
		cmd.UserFlags.AnsibleDryRun.DefaultBool,
		cmd.UserFlags.AnsibleDryRun.Usage,
	)
	flag.Parse()

	a.checkFlags()

	cmd.UserFlags.InstanceIDs.ValueStringArr = append(cmd.UserFlags.InstanceIDs.ValueStringArr,
		strings.Split(a.buffInstanceFlag, ",")...)

	return cmd.UserFlags
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.ICmd {
	a.logger = logger.Named("cmd")

	return a
}

func (a Adapter) isAllowedMode() bool {
	for _, mode := range cmd.UserFlags.Mode.AllowedValuesStr {
		if *cmd.UserFlags.Mode.ValueString == mode {
			return true
		}
	}

	return false
}

func (a *Adapter) checkFlags() {
	// check if Instance ID is defined
	if a.buffInstanceFlag == "" {
		a.logger.Error("instance IDs not defined")
		flag.PrintDefaults()

		os.Exit(1)
	}

	//  check if modes are allowed
	if !a.isAllowedMode() {
		a.logger.Error("only bash script and ansible playbook modes types are supported")
		flag.PrintDefaults()

		os.Exit(1)
	}

	// check if ansible playbook is defined
	if *cmd.UserFlags.Mode.ValueString == "ansible" &&
		*cmd.UserFlags.AnsiblePlaybook.ValueString == "" {
		a.logger.Error("running in Ansible mode but no Ansible Playbook file defined!")
		flag.PrintDefaults()

		os.Exit(1)
	}
}

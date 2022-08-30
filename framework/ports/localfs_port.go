package ports

import (
	"github.com/Trapesys/aws-commander/framework/adapters/types/ssm"
	"github.com/hashicorp/go-hclog"
)

type ILocalFS interface {
	// ReadBashScript is a file reader intended for reading bash scripts.
	ReadBashScript(bashScriptLocation string) string
	// ReadAnsiblePlaybook is a file reader intended for reading ansible playbook yaml/yml files
	ReadAnsiblePlaybook(playbookLocation string) string
	// WithLogger injects a logger instance
	WithLogger(logger hclog.Logger) ILocalFS
	// WriteRunCommandOutput writes output result from ssm.RunCommand
	WriteRunCommandOutput(cmdOutput ssm.Instances, outputLocation string) error
}

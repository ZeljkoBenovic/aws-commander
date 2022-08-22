package ports

import (
	"github.com/Trapesys/aws-commander/framework/adapters/types/ssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
)

type ISSMPort interface {
	// RunCommand runs the commands
	RunCommand() ssm.Instances

	// GetCommandOutput outputs the command result
	GetCommandOutput(commandID, instanceID *string) (string, string)

	// WithLogger builder builds the logger instance
	WithLogger(logger hclog.Logger) ISSMPort

	// WithCommands takes in the commands to execute on the instance.
	//
	// Commands are split by new line character
	WithCommands(commandsString string) ISSMPort

	// WithFreeFormCommand takes in a single bash command to be executed
	WithFreeFormCommand(cmd string) ISSMPort

	// WithInstances builder sets the instance IDs to run the command on
	//
	// Can be omitted if WithInstanceTags is set
	WithInstances(instanceIDs []string) ISSMPort

	// WithInstanceTags builder sets the instance tags to run the command on
	//
	// Can be omitted if WithInstances is set
	WithInstanceTags(tagName string, tagValues []string) ISSMPort

	// WithAWSSession builder builds the SSM adapter with AWS session
	WithAWSSession(awsSession *session.Session) ISSMPort
}

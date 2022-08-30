package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type commandType string

var generateCommand map[commandType]func() *ssm.SendCommandInput

// these should reflect flag options for -mode flag - bash by default
const (
	bashScript      commandType = "bash"
	ansiblePlaybook commandType = "ansible"
)

// prepareCommand sets the function which initializes command based on the flag input
func (a *Adapter) prepareCommand() map[commandType]func() *ssm.SendCommandInput {
	generateCommand = map[commandType]func() *ssm.SendCommandInput{
		bashScript:      a.runBashScript,
		ansiblePlaybook: a.runAnsiblePlaybook,
	}

	return generateCommand
}

// runBashScript initializes AWS-RunShellScript document which will run a bash script on nodes
func (a *Adapter) runBashScript() *ssm.SendCommandInput {
	// if we have instance IDs, else if we have tags
	if a.instanceIDs != nil {
		return &ssm.SendCommandInput{
			DocumentName:    aws.String("AWS-RunShellScript"),
			DocumentVersion: aws.String("$LATEST"),
			InstanceIds:     a.instanceIDs,
			Parameters:      a.commands,
			TimeoutSeconds:  aws.Int64(300),
		}
	} else if a.targets != nil {
		return &ssm.SendCommandInput{
			DocumentName:    aws.String("AWS-RunShellScript"),
			DocumentVersion: aws.String("$LATEST"),
			Targets:         a.targets,
			Parameters:      a.commands,
			TimeoutSeconds:  aws.Int64(300),
		}
	}

	return nil
}

func (a *Adapter) runAnsiblePlaybook() *ssm.SendCommandInput {
	// if we have instance IDs, else if we have tags
	if a.instanceIDs != nil {
		return &ssm.SendCommandInput{
			DocumentName:    aws.String("AWS-RunAnsiblePlaybook"),
			DocumentVersion: aws.String("$LATEST"),
			InstanceIds:     a.instanceIDs,
			Parameters:      a.commands,
			TimeoutSeconds:  aws.Int64(300),
		}
	} else if a.targets != nil {
		return &ssm.SendCommandInput{
			DocumentName:    aws.String("AWS-RunAnsiblePlaybook"),
			DocumentVersion: aws.String("$LATEST"),
			Targets:         a.targets,
			Parameters:      a.commands,
			TimeoutSeconds:  aws.Int64(300),
		}
	}

	return nil
}

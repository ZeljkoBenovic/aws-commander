package ssm

import (
	"fmt"
	ssm2 "github.com/Trapesys/aws-commander/framework/adapters/types/ssm"
	"github.com/Trapesys/aws-commander/framework/ports"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hashicorp/go-hclog"
	"strings"
	"time"
)

// Adapter is the adapter for SSM port
type Adapter struct {
	ssmSession          *ssm.SSM
	logger              hclog.Logger
	instanceIDs         []*string
	targets             []*ssm.Target
	mode                commandType
	ansiblePlaybookOpts *ports.AnsiblePlaybookOpts

	commands
}

type commands map[string][]*string

// NewAdapter returns new ssm instance
func NewAdapter() ports.ISSMPort {
	return &Adapter{
		commands:    make(map[string][]*string),
		instanceIDs: nil,
		targets:     nil,
	}
}

func newInstances(instances []string) []*string {
	inst := []*string{}

	for _, ins := range instances {
		inst = append(inst, aws.String(ins))
	}

	return inst
}

func (a *Adapter) WithAWSSession(awsSession *session.Session) ports.ISSMPort {
	a.ssmSession = ssm.New(awsSession)

	return a
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.ISSMPort {
	a.logger = logger.Named("ssm")

	return a
}

func (a *Adapter) WithInstances(instanceIDs []string) ports.ISSMPort {
	a.instanceIDs = newInstances(instanceIDs)

	return a
}

func (a *Adapter) WithInstanceTags(tagName string, tagValues []string) ports.ISSMPort {
	val := []*string{}
	for _, tag := range tagValues {
		val = append(val, aws.String(tag))
	}

	tgName := "tag:" + tagName

	a.targets = append(a.targets, &ssm.Target{
		Key:    &tgName,
		Values: val,
	})

	a.logger.Debug("target tags set", "target", fmt.Sprintf("%+v\n", a.targets))

	// TODO: remove panic once resolved
	panic("this feature does not work properly at this moment")
	//nolint
	return a
}

func (a *Adapter) RunCommand() ssm2.Instances {
	var command *ssm.SendCommandInput

	responseData := ssm2.Instances{}
	prepareCommandFunc := a.prepareCommand()

	if command = prepareCommandFunc[a.mode](); command == nil {
		a.logger.Error("could not find instance or tag to run the commands on")

		return responseData
	}
	// Send Command
	out, err := a.ssmSession.SendCommand(command)
	if err != nil {
		a.logger.Error("could not run SSM command", "err", err.Error())

		return responseData
	}

	a.logger.Debug("send command call response", "response", out.String())

	// TODO: when using with tags, we do not have any instance ID
	// TODO: we need some kind of mechanism to get instance ids from tags or not use Tag at all

	// wait for it to complete before returning
	for _, instance := range out.Command.InstanceIds {
		if cmdTimeoutErr := a.ssmSession.WaitUntilCommandExecutedWithContext(
			aws.BackgroundContext(),
			&ssm.GetCommandInvocationInput{
				CommandId:  out.Command.CommandId,
				InstanceId: instance,
				PluginName: nil,
			}, func(waiter *request.Waiter) {
				waiter.Delay = request.ConstantWaiterDelay(10 * time.Second)
				waiter.MaxAttempts = 60
			},
		); cmdTimeoutErr != nil {
			//nolint
			if awsErr, ok := cmdTimeoutErr.(awserr.Error); ok {
				if awsErr.Code() == request.WaiterResourceNotReadyErrorCode {
					a.logger.Error(
						"error running ansible playbook",
						"instanceID", *instance,
						"commandID", *out.Command.CommandId,
					)
				} else if awsErr.Code() == request.ErrCodeResponseTimeout {
					a.logger.Error(
						"timeout reached while waiting for command to finish",
						"instanceID", *instance,
						"commandID", *out.Command.CommandId,
					)
				}
			}
		}
	}
	// and parse and return output data
	for _, instanceID := range out.Command.InstanceIds {
		var output, errorOutput string
		output, errorOutput = a.GetCommandOutput(out.Command.CommandId, instanceID)
		inst := ssm2.Instance{
			ID:            *instanceID,
			CommandOutput: output,
			ErrorOutput:   errorOutput,
		}

		responseData.Instance = append(responseData.Instance, inst)
	}

	return responseData
}

func (a *Adapter) GetCommandOutput(cmdID, instanceID *string) (string, string) {
	out, err := a.ssmSession.GetCommandInvocation(&ssm.GetCommandInvocationInput{
		CommandId:  cmdID,
		InstanceId: instanceID,
		PluginName: nil,
	})
	if err != nil {
		a.logger.Error("could not get command output", "cmd_id", cmdID, "err", err.Error())

		return "", ""
	}

	return *out.StandardOutputContent, *out.StandardErrorContent
}

func (a *Adapter) WithCommands(commandsString string) ports.ISSMPort {
	for _, cmd := range strings.Split(commandsString, "\n") {
		a.commands["commands"] = append(a.commands["commands"], aws.String(cmd))
	}

	for _, cmdPtr := range a.commands["commands"] {
		a.logger.Debug("parsed commands to run", "commands", *cmdPtr)
	}

	return a
}

func (a *Adapter) WithFreeFormCommand(cmd string) ports.ISSMPort {
	shell := "#!/bin/bash"
	a.commands["commands"] = append(a.commands["commands"], &shell)
	a.commands["commands"] = append(a.commands["commands"], &cmd)

	return a
}

func (a *Adapter) WithMode(mode string) {
	a.mode = commandType(mode)
}

func (a *Adapter) WithAnsiblePlaybook(opts *ports.AnsiblePlaybookOpts) ports.ISSMPort {
	a.ansiblePlaybookOpts = opts

	a.commands["playbook"] = []*string{&opts.Playbook}
	a.commands["playbookurl"] = []*string{&opts.PlaybookURL}
	a.commands["extravars"] = []*string{&opts.ExtraVars}
	a.commands["check"] = []*string{&opts.Check}

	return a
}

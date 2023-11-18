package ssm

import (
	"bytes"
	"fmt"
	"github.com/Trapesys/aws-commander/conf"
	"github.com/Trapesys/aws-commander/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/davecgh/go-spew/spew"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	assm "github.com/aws/aws-sdk-go/service/ssm"
)

type ssm struct {
	log  logger.Logger
	conf conf.Config

	cl *assm.SSM
}

func (s ssm) RunBash() error {
	s.log.Info("Running ssm bash command")

	command, err := s.cl.SendCommand(&assm.SendCommandInput{
		DocumentName:    aws.String("AWS-RunShellScript"),
		DocumentVersion: aws.String("$LATEST"),
		InstanceIds:     s.provideInstanceIDs(),
		Parameters:      s.provideBashCommands(),
		TimeoutSeconds:  aws.Int64(300),
	})
	if err != nil {
		return err
	}

	s.log.Info("Command deployed successfully")
	s.log.Info("Waiting for results")

	var instIdsSuccess = make([]*string, 0)

	for _, instId := range command.Command.InstanceIds {
		if werr := s.waitForCmdExecutionComplete(command.Command.CommandId, instId); werr != nil {
			s.log.Error("Error waiting for command execution", "err", err.Error(), "instance_id", *instId)
		} else {
			instIdsSuccess = append(instIdsSuccess, instId)
		}
	}

	for _, id := range instIdsSuccess {
		out, err := s.cl.GetCommandInvocation(&assm.GetCommandInvocationInput{
			CommandId:  command.Command.CommandId,
			InstanceId: id,
		})
		if err != nil {
			s.log.Error("Could not get command output", "err", "instance_id", *id)
		} else {
			displayResults(id, out)
		}
	}

	return nil
}

func (s ssm) RunAnsible() error {
	s.log.Info("Running ssm ansible command")
	// TODO: implement
	return nil
}

func New(log logger.Logger, conf conf.Config, session *session.Session) *ssm {
	return &ssm{
		log:  log.Named("ssm"),
		conf: conf,
		cl:   assm.New(session),
	}
}

func (s ssm) provideInstanceIDs() []*string {
	var instIDs []*string

	ids := strings.Split(strings.TrimSpace(s.conf.AWSInstanceIDs), ",")
	for _, i := range ids {
		trimed := strings.TrimSpace(i)
		instIDs = append(instIDs, &trimed)
	}

	s.log.Debug("Instance ids", "ids", spew.Sdump(instIDs))

	return instIDs
}

func (s ssm) provideBashCommands() map[string][]*string {
	var (
		resp    = map[string][]*string{}
		shebang = "#!/bin/bash"
	)

	if s.conf.BashOneLiner != "" {
		resp["commands"] = append(resp["commands"], &shebang)
		resp["commands"] = append(resp["commands"], &s.conf.BashOneLiner)
	} else if s.conf.BashFile != "" {
		cmds, err := s.readBashFileAndProvideCommands()
		if err != nil {
			s.log.Fatalln("Could not provide bash commands", "err", err.Error())
		}

		for _, c := range cmds {
			resp["commands"] = append(resp["commands"], c)
		}
	} else {
		s.log.Fatalln("Bash command or bash script not specified")
	}

	s.log.Debug("Parsed commands from bash script", "cmds", spew.Sdump(resp))
	return resp
}

func (s ssm) readBashFileAndProvideCommands() ([]*string, error) {
	var cmds []*string

	fileBytes, err := os.ReadFile(s.conf.BashFile)
	if err != nil {
		return nil, err
	}

	for _, cmdLine := range strings.Split(string(fileBytes), "\n") {
		cmds = append(cmds, &cmdLine)
	}

	return cmds, nil
}

func (s ssm) waitForCmdExecutionComplete(cmdId *string, instId *string) error {
	return s.cl.WaitUntilCommandExecutedWithContext(aws.BackgroundContext(), &assm.GetCommandInvocationInput{
		CommandId:  cmdId,
		InstanceId: instId,
	}, func(waiter *request.Waiter) {
		waiter.Delay = request.ConstantWaiterDelay(time.Second * time.Duration(10))
	})
}

func displayResults(instanceId *string, data *assm.GetCommandInvocationOutput) {
	buff := bytes.Buffer{}

	buff.WriteString(fmt.Sprintf("==== INSTANCE ID - %s =====\n", *instanceId))

	if *data.StandardOutputContent != "" {
		buff.WriteString("[COMMAND OUTPUT]\n")
		buff.WriteString(*data.StandardOutputContent)
		buff.WriteString("\n")
	}

	if *data.StandardErrorContent != "" {
		buff.WriteString("[COMMAND ERROR]\n")
		buff.WriteString(*data.StandardErrorContent)
	}

	buff.WriteString("====================\n\n")

	fmt.Print(buff.String())
}

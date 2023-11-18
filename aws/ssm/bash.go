package ssm

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	assm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/davecgh/go-spew/spew"
)

func (s ssm) RunBash() error {
	s.log.Info("Running ssm bash command")

	command, err := s.cl.SendCommand(&assm.SendCommandInput{
		DocumentName:    aws.String("AWS-RunShellScript"),
		DocumentVersion: aws.String("$LATEST"),
		InstanceIds:     s.provideInstanceIDs(),
		Parameters:      s.provideBashCommands(),
		TimeoutSeconds:  &s.conf.CommandExecMaxWait,
	})
	if err != nil {
		return err
	}

	s.log.Info("Bash command deployed successfully")
	s.log.Info("Waiting for results...")

	s.waitForCmdExecAndDisplayCmdOutput(command)

	return nil
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
	var cmds = make([]*string, 0)

	fileBytes, err := os.ReadFile(s.conf.BashFile)
	if err != nil {
		return nil, err
	}

	s.log.Debug("Script content read", "content", string(fileBytes))

	for _, cmdLine := range strings.Split(string(fileBytes), "\n") {
		cmdLine := cmdLine // closure capture
		s.log.Debug("Script line read", "line", cmdLine)

		cmds = append(cmds, &cmdLine)
	}

	return cmds, nil
}

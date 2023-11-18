package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	assm "github.com/aws/aws-sdk-go/service/ssm"
	"os"
)

func (s ssm) RunAnsible() error {
	s.log.Info("Running ssm ansible command")

	command, err := s.cl.SendCommand(&assm.SendCommandInput{
		DocumentName:    aws.String("AWS-RunAnsiblePlaybook"),
		DocumentVersion: aws.String("$LATEST"),
		InstanceIds:     s.provideInstanceIDs(),
		Parameters:      s.provideAnsibleCommands(),
		TimeoutSeconds:  &s.conf.CommandExecMaxWait,
	})
	if err != nil {
		return err
	}

	s.log.Info("Command deployed successfully")
	s.log.Info("Waiting for results")
	return nil
}

func (s ssm) provideAnsibleCommands() map[string][]*string {
	var resp = map[string][]*string{}
	checkStr := "False"

	playbookStr, err := os.ReadFile(s.conf.AnsiblePlaybook)
	if err != nil {
		s.log.Fatalln("Could not read ansible playbook", "err", err.Error())
	}
	playbook := string(playbookStr)
	resp["playbook"] = []*string{&playbook}

	if s.conf.AnsibleDryRun {
		checkStr = "True"
	}
	resp["check"] = []*string{&checkStr}

	// TODO: implement "ploybookurl" and "extravars"
	resp["playbookurl"] = []*string{}
	resp["extravars"] = []*string{}

	return resp
}

package ssm

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	assm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/davecgh/go-spew/spew"
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

	s.log.Info("Ansible playbook deployed successfully")
	s.log.Info("Waiting for results...")

	s.waitForCmdExecAndDisplayCmdOutput(command)

	return nil
}

func (s ssm) provideAnsibleCommands() map[string][]*string {
	var (
		trueStr  = "True"
		falseStr = "False"
		resp     = map[string][]*string{}
		check    = map[bool]*string{
			true:  &trueStr,
			false: &falseStr,
		}
	)

	resp["check"] = []*string{check[s.conf.AnsibleDryRun]}

	if s.conf.AnsiblePlaybook != "" {
		playbookStr, err := os.ReadFile(s.conf.AnsiblePlaybook)
		if err != nil {
			s.log.Fatalln("Could not read ansible playbook", "err", err.Error())
		}

		playbook := string(playbookStr)

		resp["playbook"] = []*string{&playbook}
	}

	if s.conf.AnsibleURL != "" {
		resp["playbookurl"] = []*string{&s.conf.AnsibleURL}
	}

	if s.conf.AnsibleExtraVars != "" {
		resp["extravars"] = []*string{s.processExtraVars()}
	}

	s.log.Debug("Ansible params", "prams", spew.Sdump(resp))

	return resp
}

func (s ssm) processExtraVars() *string {
	var (
		trimmedVars   = make([]string, 0)
		processedVars string
	)

	vars := strings.Split(strings.TrimSpace(s.conf.AnsibleExtraVars), ",")
	for _, v := range vars {
		trimmedVars = append(trimmedVars, strings.TrimSpace(v))
	}

	for _, tv := range trimmedVars {
		processedVars += tv + " "
	}

	processedVars = processedVars[:len(processedVars)-1] // trim last space char

	s.log.Debug("Processed extra vars", "vars", processedVars)

	return &processedVars
}

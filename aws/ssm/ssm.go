package ssm

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Trapesys/aws-commander/conf"
	"github.com/Trapesys/aws-commander/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	assm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/davecgh/go-spew/spew"
)

type ssm struct {
	log  logger.Logger
	conf conf.Config

	cl   *assm.SSM
	sess *session.Session
}

var (
	ErrInstanceIDWaiterTimeOut = errors.New("could not get instance ids in the specified time period")
)

func New(log logger.Logger, conf conf.Config, session *session.Session) *ssm {
	return &ssm{
		log:  log.Named("ssm"),
		conf: conf,
		sess: session,
		cl:   assm.New(session),
	}
}

func (s ssm) provideTags() []*assm.Target {
	if s.conf.AWSInstanceTags == "" {
		return nil
	}

	var resp = make([]*assm.Target, 0)

	for _, rawTags := range strings.Split(s.conf.AWSInstanceTags, ";") {
		tags := make([]*string, 0)

		key := "tag:" + strings.TrimSpace(strings.Split(rawTags, "=")[0])
		vals := strings.Split(strings.Split(rawTags, "=")[1], ",")

		for _, v := range vals {
			v := strings.TrimSpace(v)
			tags = append(tags, &v)
		}

		resp = append(resp, &assm.Target{
			Key:    &key,
			Values: tags,
		})
	}

	s.log.Debug("Provided tags", "tags", spew.Sdump(resp))

	return resp
}

func (s ssm) provideInstanceIDs() []*string {
	if s.conf.AWSInstanceIDs == "" {
		return nil
	}

	var instIDs = make([]*string, 0)

	ids := strings.Split(strings.TrimSpace(s.conf.AWSInstanceIDs), ",")
	for _, i := range ids {
		trimed := strings.TrimSpace(i)
		instIDs = append(instIDs, &trimed)
	}

	s.log.Debug("Instance ids", "ids", spew.Sdump(instIDs))

	return instIDs
}

func (s ssm) waitForCmdExecutionComplete(cmdID *string, instID *string) error {
	return s.cl.WaitUntilCommandExecutedWithContext(aws.BackgroundContext(), &assm.GetCommandInvocationInput{
		CommandId:  cmdID,
		InstanceId: instID,
	}, func(waiter *request.Waiter) {
		waiter.Delay = request.ConstantWaiterDelay(time.Second * time.Duration(s.conf.CommandResultMaxWait))
	})
}

func (s ssm) waitForCmdExecAndDisplayCmdOutput(command *assm.SendCommandOutput) {
	var (
		instIdsSuccess = make([]*string, 0)
		err            error
	)
	s.log.Debug("Command", "command", spew.Sdump(command))

	instIdsSuccess, err = s.getInstanceIDsFromTags(command.Command.CommandId)
	if err != nil {
		s.log.Fatalln("Could not get instance id information", "err", err.Error())
	}

	//if ids != nil {
	//	instIds = ids
	//}
	//
	//for _, instID := range instIds {
	//	if err := s.waitForCmdExecutionComplete(command.Command.CommandId, instID); err != nil {
	//		s.log.Error("Error waiting for command execution", "err", err.Error(), "instance_id", *instID)
	//	} else {
	//		instIdsSuccess = append(instIdsSuccess, instID)
	//	}
	//}

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
}

func displayResults(instanceID *string, data *assm.GetCommandInvocationOutput) {
	buff := bytes.Buffer{}

	buff.WriteString(fmt.Sprintf("==== INSTANCE ID - %s =====\n", *instanceID))

	if *data.StandardOutputContent != "" {
		buff.WriteString("[COMMAND OUTPUT]\n")
		buff.WriteString(*data.StandardOutputContent)
	}

	if *data.StandardErrorContent != "" {
		buff.WriteString("[COMMAND ERROR]\n")
		buff.WriteString(*data.StandardErrorContent)
	}

	if *data.StandardOutputContent == "" && *data.StandardErrorContent == "" {
		buff.WriteString("NO CONTENT TO SHOW\n")
	}

	buff.WriteString("====================\n\n")

	fmt.Print(buff.String())
}

func (s ssm) getInstanceIDsFromTags(commandId *string) ([]*string, error) {
	var (
		resp    = make([]*string, 0)
		timeout = time.After(time.Duration(s.conf.CommandResultMaxWait) * time.Second)
	)

	for {
		listOut, err := s.cl.ListCommandInvocations(&assm.ListCommandInvocationsInput{
			CommandId: commandId,
		})
		if err != nil {
			return nil, err
		}

		totalInvocations := len(listOut.CommandInvocations)

		s.log.Debug("ListOut", "list", spew.Sdump(listOut))

		for _, inv := range listOut.CommandInvocations {
			if *inv.Status == assm.CommandStatusSuccess {
				resp = append(resp, inv.InstanceId)

				if len(resp) == totalInvocations {
					return resp, nil
				}
			}
		}

		select {
		case <-timeout:
			if len(resp) != 0 {
				s.log.Warn("Could not get all instance ids, returning partial results")
				s.log.Warn("For more information enable debug log")

				return resp, nil
			}

			return nil, ErrInstanceIDWaiterTimeOut
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

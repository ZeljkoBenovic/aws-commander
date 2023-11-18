package ssm

import (
	"bytes"
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

	cl *assm.SSM
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

func (s ssm) waitForCmdExecutionComplete(cmdId *string, instId *string) error {
	return s.cl.WaitUntilCommandExecutedWithContext(aws.BackgroundContext(), &assm.GetCommandInvocationInput{
		CommandId:  cmdId,
		InstanceId: instId,
	}, func(waiter *request.Waiter) {
		waiter.Delay = request.ConstantWaiterDelay(time.Second * time.Duration(s.conf.CommandResultMaxWait))
	})
}

func displayResults(instanceId *string, data *assm.GetCommandInvocationOutput) {
	buff := bytes.Buffer{}

	buff.WriteString(fmt.Sprintf("==== INSTANCE ID - %s =====\n", *instanceId))

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

	buff.WriteString("====================\n")

	fmt.Print(buff.String())
}

package ports

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
)

type CorePort interface {
	WithLogger(logger hclog.Logger) CorePort
	WithAWSSessionOptions(awsOpt session.Options) CorePort
	GetAWSSession() *session.Session
}

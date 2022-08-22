package core

import (
	"github.com/Trapesys/aws-commander/internal/ports"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
)

// Adapter plugs into Core port
type Adapter struct {
	awsSession *session.Session
	logger     hclog.Logger
}

// NewAdapter creates a new adapter
func NewAdapter() ports.CorePort {
	return &Adapter{}
}

// WithLogger builder builds a new logger instance
func (a *Adapter) WithLogger(logger hclog.Logger) ports.CorePort {
	a.logger = logger.Named("core")

	return a
}

// WithAWSSessionOptions builder builds a new AWS session instance
func (a *Adapter) WithAWSSessionOptions(awsOpt session.Options) ports.CorePort {
	var err error

	a.awsSession, err = session.NewSessionWithOptions(awsOpt)
	if err != nil {
		a.logger.Error("could not create new AWS session", "err", err.Error())

		return nil
	}

	return a
}

func (a *Adapter) GetAWSSession() *session.Session {
	return a.awsSession
}

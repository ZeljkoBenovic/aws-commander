package aws

import (
	"github.com/Trapesys/aws-commander/aws/ssm"
	"github.com/Trapesys/aws-commander/conf"
	"github.com/Trapesys/aws-commander/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

type mode string

const (
	bash    mode = "bash"
	ansible mode = "ansible"
)

type modeHandler func() error

type modesFactory map[mode]modeHandler

var (
	ErrModeNotSupported = errors.New("selected mode not supported")
)

type SSM interface {
	RunBash() error
	RunAnsible() error
}

type Aws struct {
	conf  conf.Config
	ssm   SSM
	modes modesFactory
}

func New(conf conf.Config, log logger.Logger) Aws {
	sess, err := provideSesson(conf)
	if err != nil {
		log.Fatalln("Could not create AWS session", "err", err.Error())
	}

	localssm := ssm.New(log, conf, sess)

	return Aws{
		conf: conf,
		ssm:  localssm,
		modes: modesFactory{
			bash:    localssm.RunBash,
			ansible: localssm.RunAnsible,
		},
	}
}

func (a *Aws) Run() error {
	modeHn, ok := a.modes[mode(a.conf.Mode)]
	if !ok {
		return ErrModeNotSupported
	}

	return modeHn()
}

func provideSesson(conf conf.Config) (*session.Session, error) {
	sessOpt := session.Options{}
	sessConf := aws.Config{}

	if conf.AWSRegion != "" {
		sessConf.Region = &conf.AWSRegion
	}

	if conf.AWSProfile != "" {
		sessOpt.Profile = conf.AWSProfile
	}

	sessOpt.Config = sessConf
	sessOpt.SharedConfigState = session.SharedConfigEnable

	return session.NewSessionWithOptions(sessOpt)
}

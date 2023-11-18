package app

import (
	"os"

	"github.com/Trapesys/aws-commander/aws"
	"github.com/Trapesys/aws-commander/conf"
	"github.com/Trapesys/aws-commander/logger"
	"go.uber.org/fx"
)

func Run() {
	fx.New(
		fx.Provide(
			conf.New,
			logger.New,
			aws.New,
		),
		fx.Invoke(mainApp),
		fx.NopLogger,
	).Run()
}

func mainApp(log logger.Logger, awss aws.Aws) {
	if err := awss.Run(); err != nil {
		log.Error("Run command error", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}

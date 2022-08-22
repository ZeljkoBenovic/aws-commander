package main

import (
	"github.com/ZeljkoBenovic/aws-commander/framework/adapters/left/cmd"
	"github.com/ZeljkoBenovic/aws-commander/framework/adapters/left/localfs"
	"github.com/ZeljkoBenovic/aws-commander/framework/adapters/right/ssm"
	"github.com/ZeljkoBenovic/aws-commander/internal/adapters/app"
	"github.com/ZeljkoBenovic/aws-commander/internal/adapters/core"
	"github.com/hashicorp/go-hclog"
	"os"
)

func main() {
	// init logger instance
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "aws-commander",
		Level: hclog.NoLevel,
	})

	// inject adapters into App
	commander := app.NewAdapter(
		core.NewAdapter(),
		localfs.NewAdapter(),
		ssm.NewAdapter(),
		cmd.NewAdapter(),
	).WithLogger(logger).Init()

	// run command and check for error
	if err := commander.RunCommand(); err != nil {
		logger.Error("could not run command: ", "err", err.Error())
		os.Exit(1)
	}
}

package localfs

import (
	"encoding/json"
	"fmt"
	"github.com/Trapesys/aws-commander/framework/adapters/types/ssm"
	"github.com/Trapesys/aws-commander/framework/ports"
	"github.com/hashicorp/go-hclog"
	"os"
)

type Adapter struct {
	logger hclog.Logger
}

func NewAdapter() ports.ILocalFS {
	return &Adapter{}
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.ILocalFS {
	a.logger = logger.Named("localfs")

	return a
}

func (a Adapter) ReadBashScript(bashScriptLocation string) string {
	fileBytes, err := os.ReadFile(bashScriptLocation)
	if err != nil {
		a.logger.Error("could not read file", "file", bashScriptLocation, "err", err.Error())
		os.Exit(1)
	}

	return string(fileBytes)
}

func (a Adapter) WriteRunCommandOutput(cmdOutput ssm.Instances, outputLocation string) error {
	jsonBuff, err := json.MarshalIndent(cmdOutput, "", "    ")
	if err != nil {
		return fmt.Errorf("could not marshal command output to json: %w", err)
	}

	if wrErr := os.WriteFile(outputLocation, jsonBuff, 0600); wrErr != nil {
		return fmt.Errorf("could not write file to disk: %w", wrErr)
	}

	return nil
}

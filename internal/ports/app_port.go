package ports

import (
	"github.com/hashicorp/go-hclog"
)

type IApp interface {
	RunCommand() error
	Init() IApp
	WithLogger(logger hclog.Logger) IApp
}

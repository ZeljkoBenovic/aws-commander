package ports

import (
	"github.com/ZeljkoBenovic/aws-commander/framework/adapters/types/cmd"
	"github.com/hashicorp/go-hclog"
)

type ICmd interface {
	WithLogger(logger hclog.Logger) ICmd
	GetFlags() cmd.Flags
}

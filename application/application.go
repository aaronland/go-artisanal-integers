package application

import (
	"context"
	"flag"
)

type Application interface {
	Run(context.Context, *flag.FlagSet) error
}

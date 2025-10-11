package app

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
)

type CliOption interface {
	Flags() cli.NamedFlagSets
	Validate() field.ErrorList
}

type CompleteOption interface {
	Complete() error
}

type PrintOption interface {
	String() string
}

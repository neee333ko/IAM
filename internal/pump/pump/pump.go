package pump

import (
	"context"
	"errors"

	"github.com/neee333ko/IAM/internal/pump/analytics"
)

type Pump interface {
	GetName() string
	New() Pump
	Init(interface{}) error
	WriteData(context.Context, []interface{}) error
	SetFilters(*analytics.AnalyticsFilters)
	GetFilters() *analytics.AnalyticsFilters
	SetTimeout(timeout int)
	GetTimeout() int
	SetOmitDetailedRecording(bool)
	GetOmitDetailedRecording() bool
}

func GetPumpByName(name string) (Pump, error) {
	if pump, ok := availablePumps[name]; ok && pump != nil {
		return pump, nil
	}

	return nil, errors.New(name + " Not found")
}

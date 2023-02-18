package checks

import (
	"github.com/sleepycrew/appmonitor-client/pkg/check"
	"github.com/sleepycrew/appmonitor-client/pkg/data/result"
)

type StaticCheck struct {
	Result result.Code
	Value  string
}

func (c StaticCheck) RunCheck(results chan<- check.Result) {
	results <- check.Result{
		Value:  c.Value,
		Result: c.Result,
	}
}

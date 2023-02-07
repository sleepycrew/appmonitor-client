package checks

import (
	. "github.com/sleepycrew/appmonitor-client/pkg/data"
	. "github.com/sleepycrew/appmonitor-client/pkg/data/result"
)

type StaticCheck struct {
	Name  string
	Value string
}

func (c StaticCheck) GetName() string {
	return c.Name
}

func (c StaticCheck) GetDescription() *string {
	return nil
}

func (c StaticCheck) RunCheck(results chan<- ClientCheck) {
	results <- ClientCheck{
		Name:        c.GetName(),
		Value:       c.Value,
		Description: "This is a static check for testing.",
		Result:      int(OK),
	}
}

type BadStaticCheck struct {
	Name  string
	Value string
}

func (c BadStaticCheck) GetName() string {
	return c.Name
}

func (c BadStaticCheck) GetDescription() *string {
	return nil
}

func (c BadStaticCheck) RunCheck(results chan<- ClientCheck) {
	results <- ClientCheck{
		Name:   c.GetName(),
		Value:  c.Value,
		Result: int(Error),
	}
}

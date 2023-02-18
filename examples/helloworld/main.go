package main

import (
	. "github.com/sleepycrew/appmonitor-client"
	"github.com/sleepycrew/appmonitor-client/checks"
	"github.com/sleepycrew/appmonitor-client/pkg/check"
	"github.com/sleepycrew/appmonitor-client/pkg/data/result"
	"github.com/sleepycrew/appmonitor-client/pkg/monitor"
)

func main() {
	myMonitor := NewMonitor(monitor.MonitorMetadata{
		Host:    "Hello Host",
		Website: "http://127.0.0.1",
		Ttl:     20,
	})
	var check1Name = "Check1"
	staticCheck := checks.StaticCheck{Value: "Hewwo Wowld", Result: result.OK}
	var check3Name = "Check 2"
	staticCheck3 := checks.StaticCheck{Value: "Hewwo Wowld", Result: result.OK}
	var badCheckName = "Check 3"
	badCheck := checks.StaticCheck{Value: "Hewwo Wowld", Result: result.Unknown}
	myMonitor.AddCheck(check.Metadata{Name: check1Name}, staticCheck)
	myMonitor.AddCheck(check.Metadata{Name: check3Name}, staticCheck3)
	myMonitor.AddNestedCheck(&check3Name, check.Metadata{Name: "Check 4"}, checks.StaticCheck{Value: "Hewwo Wowld"})
	myMonitor.AddNestedCheck(&check3Name, check.Metadata{Name: badCheckName}, badCheck)
	myMonitor.AddNestedCheck(&badCheckName, check.Metadata{
		Name:        "check 5",
		Description: "check if ssh exists",
	}, checks.NewFileCheck("/usr/bin/ssh", checks.FileCheckSettings{
		Exists: true,
	}))
	myMonitor.AddNestedCheck(&check3Name, check.Metadata{
		Name: "check 333",
	}, checks.StaticCheck{Value: "Hewwo Wowld"})

	var check6Name = "check 6"
	check6 := checks.StaticCheck{Value: "Hewwo Wowld", Result: result.OK}
	myMonitor.AddNestedCheck(&badCheckName, check.Metadata{Name: check6Name}, check6)
	myMonitor.AddNestedCheck(&check6Name, check.Metadata{Name: "check 7"}, checks.StaticCheck{Value: "Hewwo Wowld"})
	myMonitor.RunCmd()
}

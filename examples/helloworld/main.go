package main

import (
	. "github.com/sleepycrew/appmonitor-client"
	"github.com/sleepycrew/appmonitor-client/checks"
	"github.com/sleepycrew/appmonitor-client/pkg/monitor"
)

func main() {
	myMonitor := NewMonitor(monitor.MonitorMetadata{
		Host:    "Hello Host",
		Website: "http://127.0.0.1",
		Ttl:     20,
	})
	staticCheck := checks.StaticCheck{Name: "Check 1", Value: "Hewwo Wowld"}
	staticCheck3 := checks.StaticCheck{Name: "Check 2", Value: "Hewwo Wowld"}
	badCheck := checks.StaticCheck{Name: "Check 3", Value: "Hewwo Wowld"}
	myMonitor.AddCheck(staticCheck)
	myMonitor.AddCheck(staticCheck3)
	myMonitor.AddNestedCheck(&staticCheck3.Name, checks.StaticCheck{Name: "check 4", Value: "Hewwo Wowld"})
	myMonitor.AddNestedCheck(&staticCheck3.Name, badCheck)
	myMonitor.AddNestedCheck(&badCheck.Name, checks.NewFileCheck("check 5", "check if ssh exists", "/usr/bin/ssh", checks.FileCheckSettings{
		Exists: true,
	}))
	myMonitor.AddNestedCheck(&staticCheck3.Name, checks.StaticCheck{Name: "check 333", Value: "Hewwo Wowld"})

	check6 := checks.StaticCheck{Name: "check 6", Value: "Hewwo Wowld"}
	myMonitor.AddNestedCheck(&badCheck.Name, check6)
	myMonitor.AddNestedCheck(&check6.Name, checks.StaticCheck{Name: "check 7", Value: "Hewwo Wowld"})
	myMonitor.RunCmd()
}

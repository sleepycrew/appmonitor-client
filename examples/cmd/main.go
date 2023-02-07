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
	myCheck := checks.StaticCheck{Name: "StaticCheck", Value: "Hewwo Wowld"}
	myMonitor.AddCheck(myCheck)
	myMonitor.RunCmd()
}

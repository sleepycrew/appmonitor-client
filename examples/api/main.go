package main

import (
	"github.com/sleepycrew/appmonitor-client/checks"
	"github.com/sleepycrew/appmonitor-client/pkg/monitor"
)

func main() {
	myMonitor := monitor.NewMonitor(monitor.MonitorMetadata{
		Host:    "Hello Host",
		Website: "http://127.0.0.1",
		Ttl:     20,
	})
	myCheck := checks.StaticCheck{Value: "Hewwo Wowld"}
	myMonitor.AddCheck(myCheck)
	resp := myMonitor.Execute()
	println("Check Result: ", resp.Meta.Result)
}

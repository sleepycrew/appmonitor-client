package cmd

import (
	json2 "encoding/json"
	"flag"
	"github.com/sleepycrew/appmonitor-client/internal/monitor"
	"os"
)

var startServer bool
var jsonOutput bool
var help bool

func ExecuteCmd(monitor monitor.MonitorInterface) {
	flag.BoolVar(&startServer, "server", false, "prints this message")
	flag.BoolVar(&jsonOutput, "json", false, "outputs the result as json")
	flag.BoolVar(&help, "help", false, "prints this message")

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if startServer {
		println("Starting Server")
		monitor.StartServer()
		os.Exit(0)
	} else {
		result := monitor.Execute()
		if jsonOutput {
			bytes, _ := json2.Marshal(result)
			println(string(bytes))
		} else {
			println(result.Meta.Host)
			println(result.Meta.Website)
			println(result.Meta.Time, "ms")
			println("Ran ", len(result.Checks), " Checks")
			for _, check := range result.Checks {
				println(check.Name, "\t", check.Result, "\t", check.Time)
			}
		}
		os.Exit(result.Meta.Result)
	}
}

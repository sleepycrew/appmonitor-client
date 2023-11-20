# Appmonitor Go Client

This is a client implementation for the [IML Appmonitor](https://os-docs.iml.unibe.ch/appmonitor) written in go.

The client acts as an agent that can be scraped for status reports and metrics, similar to a prometheus exporter, 
but with a major focus on checks about the underlying system. For Instance there might be a check that verifies that a 
service is running or a path is writable.

The api is currently unstable, but shouldn't change drastically.

## Design
The client is not a fixed binary with a configuration file but rather a framework/library that can be used to build an agent for a specific system. It would also be possible to incorporate it into other programs as well.

All checks implemented in this repository must be platform independent, as in "work on most *nix systems", other platforms are a nice to have.

Users can easily implement their own checks and metrics and add them to their agent.

## Usage

The following example provides a binary that accepts the following flags:
- -json
  - output json instead of text
- -server
  - start a server that can be scraped

Without any flags the binary would run the checks and print the results to stdout.

```go
package main

import (
	. "github.com/sleepycrew/appmonitor-client"
	"github.com/sleepycrew/appmonitor-client/checks"
	"github.com/sleepycrew/appmonitor-client/pkg/check"
	"github.com/sleepycrew/appmonitor-client/pkg/monitor"
)

func main() {
	myMonitor := NewMonitor(monitor.MonitorMetadata{
		Host:    "Hello Host",
		Website: "http://127.0.0.1",
		Ttl:     20,
	})
	myCheck := checks.StaticCheck{Value: "Hewwo Wowld"}
	myMonitor.AddCheck(check.Metadata{Name: "StaticCheck"}, myCheck)
	myMonitor.RunCmd()
}
```

You could also change the last line to always start the server should you so desire.

```go
myMonitor.StartServer()
```

Remember `appmonitor-client` is really just a framework for building a client tailored to your needs.


## Todo
- [ ] Implement counter support
- [ ] Add auth support
- [ ] Add notification support 
- [ ] Improve logging
- [ ] Consider i18n
- [ ] Add more checks

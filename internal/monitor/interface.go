package monitor

import "github.com/sleepycrew/appmonitor-client/pkg/data"

type MonitorInterface interface {
	GetName() string
	Execute() data.ClientResponse
	StartServer()
}

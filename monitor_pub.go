package appmonitor_client

import (
	"github.com/google/uuid"
	"github.com/sleepycrew/appmonitor-client/pkg/monitor"
)

func NewMonitor(meta monitor.MonitorMetadata) monitor.Monitor {
	u := uuid.New()
	return NewNamedMonitor(u.String(), meta)
}

func NewNamedMonitor(name string, meta monitor.MonitorMetadata) monitor.Monitor {
	m := monitor.NewNamedMonitor(name, meta)
	return m
}

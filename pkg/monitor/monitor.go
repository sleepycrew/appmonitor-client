package monitor

import (
	"github.com/google/uuid"
	"github.com/sleepycrew/appmonitor-client/internal/cmd"
	"github.com/sleepycrew/appmonitor-client/internal/server"
	"github.com/sleepycrew/appmonitor-client/pkg/check"
	"github.com/sleepycrew/appmonitor-client/pkg/data"
	"time"
)

type MonitorMetadata struct {
	Host    string
	Tags    []string
	Ttl     int
	Website string
}

type Monitor struct {
	name       string
	metadata   MonitorMetadata
	checkSuite check.Checksuite
}

func (m Monitor) Execute() data.ClientResponse {
	start := time.Now()
	checks, result := m.checkSuite.RunChecks()
	elapsed := time.Since(start).Milliseconds()
	var pointyPointers []*data.ClientCheck
	for i, _ := range checks {
		p := checks[i]
		pointyPointers = append(pointyPointers, &p)
	}
	return data.ClientResponse{
		Meta: &data.ClientMeta{
			Host:    m.metadata.Host,
			Website: m.metadata.Website,
			Ttl:     m.metadata.Ttl,
			Time:    float64(elapsed),
			Result:  int(result),
		},
		Checks: pointyPointers,
	}
}

func (m Monitor) GetMetadata() MonitorMetadata {
	return m.metadata
}

func (m Monitor) StartServer() {
	server.StartServer(m)
}

func (m Monitor) RunCmd() {
	cmd.ExecuteCmd(m)
}

func (m Monitor) GetName() string {
	return m.name
}

func (m *Monitor) AddCheck(c check.Check) {
	m.checkSuite.AddCheck(c)
}

func (m *Monitor) AddNestedCheck(parent *string, c check.Check) {
	m.checkSuite.AddNestedCheck(parent, c)
}

func NewMonitor(meta MonitorMetadata) Monitor {
	u := uuid.New()
	return NewNamedMonitor(u.String(), meta)
}

func NewNamedMonitor(name string, meta MonitorMetadata) Monitor {
	m := Monitor{name, meta, check.NewCheckTreeSuite()}
	return m
}

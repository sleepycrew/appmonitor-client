package monitor

// TODO implement

// It is supposed to provide a way to orchestrate multiple monitor instances for servers that run multiple applications

type Combination struct {
	monitors map[string]Monitor
}

func (c Combination) RegisterMonitor(monitor Monitor) {
	c.monitors[monitor.GetName()] = monitor
}

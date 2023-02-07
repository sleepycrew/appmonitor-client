package check

import (
	"sync"
	"time"

	"github.com/sleepycrew/appmonitor-client/pkg/data"
)

type Check interface {
	GetName() string
	GetDescription() *string
	RunCheck(result chan<- data.ClientCheck)
}

// Runs a check and sets the Time field based on execution time
// creating a channel is probably costly?
func collectRuntime(check Check, result chan<- data.ClientCheck) {
	c := make(chan data.ClientCheck)
	start := time.Now()
	go check.RunCheck(c)
	tempRes := <-c
	elapsed := time.Since(start).Milliseconds()
	tempRes.Time = float64(elapsed)
	result <- tempRes
}

type clientCheckResolver = func(result chan<- data.ClientCheck)

func setParent(parent *string, fun clientCheckResolver, result chan<- data.ClientCheck) {
	if parent != nil {
		c := make(chan data.ClientCheck)
		go fun(c)
		tempRes := <-c
		tempRes.Parent = *parent
		result <- tempRes
	} else {
		go fun(result)
	}
}

func merge(cs ...<-chan data.ClientCheck) <-chan data.ClientCheck {
	var wg sync.WaitGroup
	out := make(chan data.ClientCheck)

	output := func(c <-chan data.ClientCheck) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

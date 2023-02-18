package check

import (
	"github.com/sleepycrew/appmonitor-client/pkg/data/result"
	"sync"
	"time"

	"github.com/sleepycrew/appmonitor-client/pkg/data"
)

type Check interface {
	RunCheck(result chan<- Result)
}

type Result struct {
	Result result.Code
	Value  string
}

type Metadata struct {
	Name        string
	Description string
}

// Runs a check and sets the Time field based on execution time
// creating a channel is probably costly?
func collectRuntime(check Check, result chan<- data.ClientCheck) {
	c := make(chan Result)
	start := time.Now()
	go check.RunCheck(c)
	checkResult := <-c
	elapsed := time.Since(start).Milliseconds()

	clientCheck := new(data.ClientCheck)
	clientCheck.Value = checkResult.Value
	clientCheck.Result = int(checkResult.Result)
	clientCheck.Time = float64(elapsed)
	result <- *clientCheck
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

func setMetadata(cc *data.ClientCheck, metadata Metadata) {
	cc.Name = metadata.Name
	cc.Description = metadata.Description
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

// Parent
// converts a string to a pointer to make AddCheck calls bearable
// why is go like this ? :(
func Parent(s string) *string {
	return &s
}

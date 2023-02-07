package check

import (
	. "github.com/sleepycrew/appmonitor-client/pkg/data"
	. "github.com/sleepycrew/appmonitor-client/pkg/data/result"
	"time"
)

type Checksuite interface {
	RunChecks() ([]ClientCheck, Result)
	AddCheck(check Check)
	AddNestedCheck(parent *string, check Check)
}

// ParallelChecksuite Runs all checks in parallel regardless of hierarchy
type ParallelChecksuite struct {
	checks map[string]Check
}

func (s ParallelChecksuite) RunChecks() ([]ClientCheck, Result) {
	c := make(chan ClientCheck)
	count := 0

	for k := range s.checks {
		check := s.checks[k]
		go collectRuntime(check, c)
		count++
	}

	results := make([]ClientCheck, count)
	for i := range results {
		results[i] = <-c
	}
	return results, Unknown
}

func (s ParallelChecksuite) AddCheck(check Check) {
	name := check.GetName()
	s.checks[name] = check
}

func (s ParallelChecksuite) AddNestedCheck(_ *string, check Check) {
	s.AddCheck(check)
}

type TreeChecksuite struct {
	tree *checkTree
}

func (s TreeChecksuite) AddCheck(check Check) {
	s.tree.AddCheck(nil, check)
}

func (s TreeChecksuite) AddNestedCheck(parent *string, check Check) {
	s.tree.AddCheck(parent, check)
}

func (s TreeChecksuite) RunChecks() ([]ClientCheck, Result) {
	size := s.tree.Size()
	println("tree size: ", size)

	channels := make([]chan ClientCheck, 0)
	for i, node := range s.tree.children {
		channels = append(channels, make(chan ClientCheck, 0))
		go evaluateTree(node, true, nil, channels[i])
	}

	results := make([]ClientCheck, 0)
	for _, channel := range channels {
		for check := range channel {
			results = append(results, check)
			//println(check.Name)
		}
	}

	println("hey2")
	return results, Unknown
}

func evaluateTree(node *checkTreeNode, parentSuccess bool, parentName *string, result chan<- ClientCheck) {
	time.Sleep(2 * time.Second)
	println(node.Value.GetName())
	defer close(result)
	if !parentSuccess {
		result <- ClientCheck{
			// handle description
			Name: node.Value.GetName(),
			Time: 0,
			// assume that parent is not nil, first call must use parentSuccess = true anyway
			Parent: *parentName,
			Result: Unknown,
			Value:  "Parent Failed.",
		}
		return
	}

	c := make(chan ClientCheck)
	go setParent(parentName, func(r chan<- ClientCheck) {
		go collectRuntime(node.Value, r)
	}, c)
	check := <-c
	close(c)

	channels := make([]<-chan ClientCheck, 0)

	for _, child := range node.Children {
		nodeName := node.Value.GetName()
		channel := make(chan ClientCheck)
		channels = append(channels, channel)
		go evaluateTree(child, check.Result != Error, &nodeName, channel)
	}

	for childCheck := range merge(channels...) {
		result <- childCheck
	}

	result <- check
}

/**
whole big mess :<

what needs to be done in easy words ->

func evaluateTree(node *checkTreeNode, parentSuccess bool, parentName *string, result chan<- ClientCheck) {
	if !parentSuccess {
		println("===========")
		println(node.Value.GetName())
		println(parentName)
		println("===========")
		result <- ClientCheck{
			// handle description
			Name:   node.Value.GetName(),
			Time:   0,
			Parent: *parentName,
			Result: Unknown,
		}
		return
	}
	childResults := make(chan ClientCheck, len(node.Children))
	go collectRuntime(node.Value, childResults)

	res := <-childResults
	for _, child := range node.Children {
		go evaluateTree(child, res.Result != Error, &res.Name, childResults)
	}
	if res.Result == Error {
		// early return nyan
		result <- res
		return
	}

	for range node.Children {
		childRes := <-childResults
		childRes.Parent = res.Name
		println("child: ", childRes.Name)
		result <- childRes
	}

	println("res: ", res.Name)
	result <- res
}


func evaluateTree(node *checkTreeNode, parentSuccess bool, result chan<- ClientCheck) {
	if !parentSuccess {
		result <- ClientCheck{
			// handle description
			Name: node.Value.GetName(),
			Time: 0,
			//Parent: *parentName,
			Result: Unknown,
		}
		return
	}
	childResults := make(chan ClientCheck, len(node.Children))
	for _, child := range node.Children {
		go evaluateTree(child, true, childResults)
	}
	go node.Value.RunCheck(childResults)
	check := <-childResults
	for range node.Children {
		var childCheck = <-childResults
		result <- childCheck
	}
	result <- check
}

*/

func NewCheckSuite() Checksuite {
	return ParallelChecksuite{checks: make(map[string]Check)}
}

func NewCheckTreeSuite() Checksuite {
	return TreeChecksuite{tree: &checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0), value: nil}}
}

package check

import (
	. "github.com/sleepycrew/appmonitor-client/pkg/data"
	"github.com/sleepycrew/appmonitor-client/pkg/data/result"
)

type Checksuite interface {
	RunChecks() ([]ClientCheck, result.Code)
	AddCheck(metadata Metadata, check Check)
	AddNestedCheck(parent *string, metadata Metadata, check Check)
}

// ParallelChecksuite Runs all checks in parallel regardless of hierarchy
type ParallelChecksuite struct {
	checks   map[string]Check
	metadata map[string]Metadata
}

func (s ParallelChecksuite) RunChecks() ([]ClientCheck, result.Code) {
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
	return results, result.Unknown
}

func (s ParallelChecksuite) AddCheck(metadata Metadata, check Check) {
	s.checks[metadata.Name] = check
	s.metadata[metadata.Name] = metadata
}

func (s ParallelChecksuite) AddNestedCheck(_ *string, metadata Metadata, check Check) {
	s.AddCheck(metadata, check)
}

type TreeChecksuite struct {
	tree *checkTree
}

func (s TreeChecksuite) AddCheck(metadata Metadata, check Check) {
	s.tree.AddCheck(nil, metadata, check)
}

func (s TreeChecksuite) AddNestedCheck(parent *string, metadata Metadata, check Check) {
	s.tree.AddCheck(parent, metadata, check)
}

func (s TreeChecksuite) RunChecks() ([]ClientCheck, result.Code) {
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
		}
	}

	return results, result.Unknown
}

func evaluateTree(node *checkTreeNode, parentSuccess bool, parentName *string, res chan<- ClientCheck) {
	defer close(res)
	if !parentSuccess {
		res <- ClientCheck{
			// handle description
			Name: node.Metadata.Name,
			Time: 0,
			// assume that parent is not nil, first call must use parentSuccess = true anyway
			Parent: *parentName,
			Result: result.Unknown,
			Value:  "Parent Failed.",
		}
		return
	}

	c := make(chan ClientCheck)
	go setParent(parentName, func(r chan<- ClientCheck) {
		go collectRuntime(node.Value, r)
	}, c)
	check := <-c
	setMetadata(&check, node.Metadata)
	close(c)

	channels := make([]<-chan ClientCheck, 0)

	for _, child := range node.Children {
		nodeName := node.Metadata.Name
		channel := make(chan ClientCheck)
		channels = append(channels, channel)
		go evaluateTree(child, check.Result != result.Error, &nodeName, channel)
	}

	for childCheck := range merge(channels...) {
		res <- childCheck
	}

	res <- check
}

func NewCheckSuite() Checksuite {
	return ParallelChecksuite{checks: make(map[string]Check)}
}

func NewCheckTreeSuite() Checksuite {
	return TreeChecksuite{tree: &checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0)}}
}

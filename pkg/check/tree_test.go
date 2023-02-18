package check

import (
	"github.com/sleepycrew/appmonitor-client/pkg/data/result"
	"github.com/stretchr/testify/assert"
	"testing"
)

// redefine (import cycle)

type StaticCheck struct {
	Result result.Code
	Value  string
}

func (c StaticCheck) RunCheck(results chan<- Result) {
	results <- Result{
		Value:  c.Value,
		Result: c.Result,
	}
}

func TestCheckTreeNode_FindByName(t *testing.T) {
	tree := checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0)}
	tree.AddCheck(nil, Metadata{
		Name: "Hello",
	},
		StaticCheck{
			Value:  "Hello World",
			Result: result.OK,
		})
	var hello2 = StaticCheck{
		Value:  "World",
		Result: result.OK,
	}
	tree.AddCheck(nil, Metadata{
		Name: "Hello2",
	}, hello2)
	tree.AddCheck(nil, Metadata{Name: "Hello4"}, StaticCheck{
		Value:  "World",
		Result: result.OK,
	})
	var hello3 = StaticCheck{
		Value: "World",
	}
	tree.AddCheck(Parent("Hello2"), Metadata{
		Name: "Hello3",
	}, hello3)
	var hello4 = StaticCheck{
		Value: "World",
	}
	tree.AddCheck(Parent("Hello3"), Metadata{
		Name: "Hello4",
	}, hello4)

	result := tree.FindByName("Hello2")
	assert.NotNil(t, result)
	assert.Equal(t, "Hello2", result.Metadata.Name)

	result2 := tree.FindByName("Hello3")
	assert.NotNil(t, result2)
	assert.Equal(t, "Hello3", result2.Metadata.Name)

	result3 := tree.FindByName("Hello4")
	assert.NotNil(t, result3)
	assert.Equal(t, "Hello4", result3.Metadata.Name)
}

func TestCheckTree_AddCheck(t *testing.T) {
	tree := checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0)}
	tree.AddCheck(nil, Metadata{
		Name: "Hello",
	}, StaticCheck{
		Value:  "World",
		Result: result.OK,
	})

	assert.Equal(t, 1, len(tree.children))
}

func TestCheckTree_AddCheckNested(t *testing.T) {
	tree := &checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0)}
	hello := StaticCheck{
		Value: "World",
	}
	tree.AddCheck(nil, Metadata{
		Name: "Hello",
	}, hello)

	tree.AddCheck(Parent("Hello"), Metadata{
		Name: "Hello2",
	}, StaticCheck{
		Value: "World",
	})

	assert.Equal(t, 1, len(tree.children[0].Children))
}

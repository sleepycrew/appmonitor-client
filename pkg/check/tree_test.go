package check

import (
	"github.com/sleepycrew/appmonitor-client/checks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckTreeNode_FindByName(t *testing.T) {
	tree := checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0), value: nil}
	tree.AddCheck(nil, checks.StaticCheck{
		Name:  "Hello",
		Value: "World",
	})
	var hello2 = checks.StaticCheck{
		Name:  "Hello2",
		Value: "World",
	}
	tree.AddCheck(nil, hello2)
	tree.AddCheck(nil, checks.StaticCheck{
		Name:  "Hello4",
		Value: "World",
	})
	var hello3 = checks.StaticCheck{
		Name:  "Hello3",
		Value: "World",
	}
	tree.AddCheck(&hello2.Name, hello3)
	var hello4 = checks.StaticCheck{
		Name:  "Hello4",
		Value: "World",
	}
	tree.AddCheck(&hello3.Name, hello4)

	result := tree.FindByName("Hello2")
	assert.NotNil(t, result)
	assert.Equal(t, "Hello2", result.Value.GetName())

	result2 := tree.FindByName("Hello3")
	assert.NotNil(t, result2)
	assert.Equal(t, "Hello3", result2.Value.GetName())

	result3 := tree.FindByName("Hello4")
	assert.NotNil(t, result3)
	assert.Equal(t, "Hello4", result3.Value.GetName())
}

func TestCheckTree_AddCheck(t *testing.T) {
	tree := checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0), value: nil}
	tree.AddCheck(nil, checks.StaticCheck{
		Name:  "Hello",
		Value: "World",
	})

	assert.Equal(t, 1, len(tree.children))
}

func TestCheckTree_AddCheckNested(t *testing.T) {
	tree := &checkTree{checkNames: map[string]bool{}, children: make([]*checkTreeNode, 0), value: nil}
	hello := checks.StaticCheck{
		Name:  "Hello",
		Value: "World",
	}
	tree.AddCheck(nil, hello)

	tree.AddCheck(&hello.Name, checks.StaticCheck{
		Name:  "Hello2",
		Value: "World",
	})

	assert.Equal(t, 1, len(tree.children[0].Children))
}

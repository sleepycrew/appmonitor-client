package check

type checkTree struct {
	// cache of all names of the checks inside the tree
	checkNames map[string]bool
	children   []*checkTreeNode
}

func (c *checkTree) FindByName(name string) *checkTreeNode {
	for _, child := range c.children {
		res := child.FindByName(name)
		if res != nil {
			return res
		}
	}

	return nil
}

func (c *checkTree) AddCheck(parent *string, metadata Metadata, check Check) int {
	checkName := metadata.Name

	if c.checkNames[checkName] {
		// duplicate
		return 1
	}

	if parent != nil {
		if !c.checkNames[*parent] {
			// parent does not exist
			return 2
		}

		p := c.FindByName(*parent)
		p.Children = append(p.Children, &checkTreeNode{
			Value:    check,
			Metadata: metadata,
			Children: make([]*checkTreeNode, 0),
		})
	} else {
		c.children = append(c.children, &checkTreeNode{
			Value:    check,
			Metadata: metadata,
			Children: make([]*checkTreeNode, 0),
		})
	}

	c.checkNames[checkName] = true
	return 0
}

func (c checkTree) Size() int {
	size := 0
	for _, node := range c.children {
		size += sizeOfTree(node)
	}
	return size
}

type checkTreeNode struct {
	Value    Check
	Metadata Metadata
	Children []*checkTreeNode
}

func (c *checkTreeNode) FindByName(name string) *checkTreeNode {
	if c.Metadata.Name == name {
		return c
	}

	for _, child := range c.Children {
		res := child.FindByName(name)
		if res != nil {
			return res
		}
	}

	return nil
}

func sizeOfTree(node *checkTreeNode) int {
	size := 1
	for _, child := range node.Children {
		size += sizeOfTree(child)
	}
	return size
}

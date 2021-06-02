package truenasapi

type (
	Node struct {
		ID         int
		Name       string
		Initiators []string
		Networks   []string
	}
)

func (c *client) GetNode(name string) (*Node, error) {
	initiator, err := c.getIntiator(name)
	if err != nil {
		return nil, err
	}

	node := &Node{
		ID:         initiator.ID,
		Name:       initiator.Comment,
		Networks:   initiator.AuthNetwork,
		Initiators: initiator.Initiators,
	}

	return node, nil
}

func (c *client) SetNode(node Node) error {
	_, err := c.GetNode(node.Name)
	if err != nil {
		if err, ok := err.(*NotFoundError); !ok {
			return err
		}

		c.DeleteNode(node.Name)
	}

	initiator := InitiatorOpts{
		Initiators:  node.Initiators,
		AuthNetwork: node.Networks,
		Comment:     node.Name,
	}

	_, err = c.createInitiator(initiator)
	return err
}

func (c *client) DeleteNode(name string) error {
	return c.deleteInitiator(name)
}

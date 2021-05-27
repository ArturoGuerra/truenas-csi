package truenasapi

type (
	Node struct {
		IP   string
		ID   int
		Name string
	}
)

func (c *client) GetNode(name string) (*Node, error) {
	return nil, nil
}

func (c *client) SetNode(node Node) error {
	return nil
}

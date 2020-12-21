package connect

var manager = &WsManager{
	clients: make(map[string]*Client),
}

func RegisterClient(c *Client) {
	manager.clients[c.getId()] = c
}

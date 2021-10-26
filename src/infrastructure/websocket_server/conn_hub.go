package websocketserver

type ConnHub struct {
	clients    map[string][]*Client
	send       chan JsonData
	register   chan *Client
	unregister chan *Client
}

func (hub *ConnHub) Run() {
	for {
		select {
		// Register client to hub
		case client := <-hub.register:
			id := client.room
			hub.clients[id] = append(hub.clients[id], client)

		// Unregister client to hub
		case client := <-hub.unregister:
			id := client.room
			if _, ok := hub.clients[id]; ok {
				delete(hub.clients, id)
				close(client.send)
			}

		// Loop through registered clients and send message to their send channel
		case message := <-hub.send:
			id := message.Room
			for range hub.clients {
				if clients, ok := hub.clients[id]; ok {
					for _, client := range clients {
						select {
						case client.send <- message:
						// If send buffer is full, assume client is dead or stuck and unregister
						default:
							close(client.send)
							delete(hub.clients, id)
						}
					}
				}
			}

		}
	}

}

func NewConnHub() *ConnHub {
	return &ConnHub{
		clients:    make(map[string][]*Client),
		send:       make(chan JsonData),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

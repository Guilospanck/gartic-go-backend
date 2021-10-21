package websocketserver

type ConnHub struct {
	clients    map[*Client]string
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (hub *ConnHub) Run() {
	for {
		select {
		// Register client to hub
		case client := <-hub.register:
			hub.clients[client] = client.username
		// Unregister client to hub
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		// Loop through registered clients and send message to their send channel
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				// If send buffer is full, assume client is dead or stuck and unregister
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}

}

func NewConnHub() *ConnHub {
	return &ConnHub{
		clients:    make(map[*Client]string),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

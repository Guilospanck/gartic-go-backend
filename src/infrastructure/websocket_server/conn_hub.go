package websocketserver

type ConnHub struct {
	clients                       map[string][]*Client
	send                          chan JsonData
	register                      chan *Client
	unregister                    chan *Client
	sendToWaitingRoom             chan []RoomAndParticipants
	broadcastToParticipantsInRoom chan RoomAndParticipants
	broadcastParticipantTurn      chan ParticipantsTurn
}

func (hub *ConnHub) removeClientFromRoomList(room string, client *Client) {
	if _, ok := hub.clients[room]; !ok {
		return
	}

	clientsFromRoom := hub.clients[room]
	if len(clientsFromRoom) == 1 {
		delete(hub.clients, room)
		return
	}

	lengthOfClientsSlice := len(clientsFromRoom)
	indexOfElement := 0
	for index, value := range clientsFromRoom {
		if value == client {
			indexOfElement = index
		}
	}

	if indexOfElement == lengthOfClientsSlice-1 { // last element
		clientsFromRoom = clientsFromRoom[:lengthOfClientsSlice-1]
	} else {
		firstSlice := clientsFromRoom[:indexOfElement]
		secondSlice := clientsFromRoom[indexOfElement+1 : lengthOfClientsSlice]
		newSlice := []*Client{}
		newSlice = append(newSlice, firstSlice...)
		newSlice = append(newSlice, secondSlice...)
		clientsFromRoom = newSlice
	}

	hub.clients[room] = clientsFromRoom
}

func (hub *ConnHub) Run() {
	for {
		select {
		// Register client to hub
		case client := <-hub.register:
			id := client.Room
			hub.clients[id] = append(hub.clients[id], client)

		// Unregister client to hub
		case client := <-hub.unregister:
			id := client.Room
			if _, ok := hub.clients[id]; ok {
				close(client.Send)
				hub.removeClientFromRoomList(client.Room, client)
			}

		// Loop through registered clients and send message to their send channel
		case message := <-hub.send:
			id := message.Room
			if clients, ok := hub.clients[id]; ok {
				for _, client := range clients {
					select {
					case client.Send <- message:
					default:
					}
				}
			}

		case message := <-hub.sendToWaitingRoom:
			id := "waitingroomgarticlikeapp"
			for range hub.clients {
				if clients, ok := hub.clients[id]; ok {
					for _, client := range clients {
						select {
						case client.Send <- message:
						default:
						}
					}
				}
			}

		case message := <-hub.broadcastToParticipantsInRoom:
			id := message.Room
			for range hub.clients {
				if clients, ok := hub.clients[id]; ok {
					for _, client := range clients {
						select {
						case client.Send <- message:
						default:
						}
					}
				}
			}

		case message := <-hub.broadcastParticipantTurn:
			id := message.Room
			if clients, ok := hub.clients[id]; ok {
				for _, client := range clients {
					select {
					case client.Send <- message:
					default:
					}
				}
			}

		}
	}

}

func NewConnHub() *ConnHub {
	return &ConnHub{
		clients:                       make(map[string][]*Client),
		send:                          make(chan JsonData),
		register:                      make(chan *Client),
		unregister:                    make(chan *Client),
		sendToWaitingRoom:             make(chan []RoomAndParticipants),
		broadcastToParticipantsInRoom: make(chan RoomAndParticipants),
		broadcastParticipantTurn:      make(chan ParticipantsTurn),
	}
}

package chat

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type ChatRoom struct {
	Hub *Hub
}

func NewChatRoom(m *ChatroomManager) *ChatRoom {
	return &ChatRoom{
		Hub: NewHub(m),
	}
}

func (cr *ChatRoom) BroadcastMessage(message Message) {
	cr.Hub.Broadcast <- message
}

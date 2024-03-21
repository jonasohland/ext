package http

type MessageBody struct {
	Message string `json:"message"`
}

func NewMessageBody(msg string) *MessageBody {
	return &MessageBody{Message: msg}
}

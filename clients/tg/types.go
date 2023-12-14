package tg

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Id      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string
	From From
	Chat Chat
}
type From struct {
	Username string `json:"username"`
}
type Chat struct {
	Id int `json:"id"`
}

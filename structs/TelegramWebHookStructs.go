package structs

type MessageBody struct {
	Text      string       `json:"text"`
	MessageId int          `json:"message_id"`
	Date      int          `json:"date"`
	From      FromBody     `json:"from"`
	Chat      FromBody     `json:"chat"`
	Location  LocationBody `json:"location"`
}

type LocationBody struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type FromBody struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type TelegramRequestBody struct {
	Message  MessageBody `json:"message"`
	UpdateId int         `json:"update_id"`
}

package request

type SendEmailRegister struct {
	UserId   string `json:"userId" bson:"userId"`
	FullName string `json:"fullName" bson:"fullName"`
	Email    string `json:"email" bson:"email"`
	Otp      string `json:"otp" bson:"otp"`
}

type EmailRequest struct {
	FullName     string `json:"fullName"`
	EventName    string `json:"eventName"`
	TicketNumber string `json:"ticketNumber"`
	OrderTime    string `json:"orderTime"`
	PaymentType  string `json:"paymentType"`
	EventTime    string `json:"eventTime"`
	EventPlace   string `json:"eventPlace"`
}

type TicketRequest struct {
	FullName     string `json:"fullName"`
	TicketType   string `json:"ticketType"`
	TicketNumber string `json:"ticketNumber"`
	TicketPrice  string `json:"ticketPrice"`
	SeatNumber   int    `json:"seatNumber"`
	EventName    string `json:"eventName"`
	EventTime    string `json:"eventTime"`
	EventPlace   string `json:"eventPlace"`
	QRCode       string `json:"qrCode"`
}

type SendTicketReq struct {
	OrderId string `json:"orderId"`
}

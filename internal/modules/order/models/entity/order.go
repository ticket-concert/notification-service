package entity

import "time"

type Country struct {
	Name  string `json:"name" bson:"name"`
	Code  string `json:"code" bson:"code"`
	City  string `json:"city" bson:"city"`
	Place string `json:"place" bson:"place"`
}

type Order struct {
	OrderId       string    `json:"orderId" bson:"orderId"`
	PaymentId     string    `json:"paymentId" bson:"paymentId"`
	MobileNumber  string    `json:"mobileNumber" bson:"mobileNumber"`
	VaNumber      string    `json:"vaNumber" bson:"vaNumber"`
	Bank          string    `json:"bank" bson:"bank"`
	Email         string    `json:"email" bson:"email"`
	FullName      string    `json:"fullName" bson:"fullName"`
	TicketNumber  string    `json:"ticketNumber" bson:"ticketNumber"`
	TicketType    string    `json:"ticketType" bson:"ticketType"`
	SeatNumber    int       `json:"seatNumber" bson:"seatNumber"`
	EventName     string    `json:"eventName" bson:"eventName"`
	Country       Country   `json:"country" bson:"country"`
	DateTime      time.Time `json:"dateTime" bson:"dateTime"`
	Description   string    `json:"description" bson:"description"`
	Tag           string    `json:"tag" bson:"tag"`
	Amount        int       `json:"amount" bson:"amount"`
	PaymentStatus string    `json:"paymentStatus" bson:"paymentStatus"`
	OrderTime     time.Time `json:"orderTime" bson:"orderTime"`
	UserId        string    `json:"userId" bson:"userId"`
	QueueId       string    `json:"queueId" bson:"queueId"`
	TicketId      string    `json:"ticketId" bson:"ticketId"`
	EventId       string    `json:"eventId" bson:"eventId"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}

package domain

type TicketBookingConfirmed struct {
	Header        EventHeader `json:"header"`
	TicketID      string      `json:"ticket_id"`
	CustomerEmail string      `json:"customer_email"`
	Price         Money       `json:"price"`
}

func (t TicketBookingConfirmed) isEvent() {
}

var _ Event = (*TicketBookingConfirmed)(nil)

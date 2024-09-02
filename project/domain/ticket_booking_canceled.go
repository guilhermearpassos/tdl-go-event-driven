package domain

type TicketBookingCanceled struct {
	Header        EventHeader `json:"header"`
	TicketID      string      `json:"ticket_id"`
	CustomerEmail string      `json:"customer_email"`
	Price         Money       `json:"price"`
}

func (t TicketBookingCanceled) isEvent() {
}

var _ Event = (*TicketBookingCanceled)(nil)

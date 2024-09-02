package domain

type TicketStatus struct {
	TicketID      string `json:"ticket_id"`
	Status        string `json:"status"`
	Price         Money  `json:"price"`
	CustomerEmail string `json:"customer_email"`
}

type TicketsStatusRequest struct {
	Tickets []TicketStatus `json:"tickets"`
}

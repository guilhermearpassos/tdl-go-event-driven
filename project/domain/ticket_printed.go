package domain

type TicketPrinted struct {
	Header   EventHeader `json:"header"`
	TicketID string      `json:"ticket_id"`
	FileName string      `json:"file_name"`
}

func (t TicketPrinted) isEvent() {
}

var _ Event = (*TicketPrinted)(nil)

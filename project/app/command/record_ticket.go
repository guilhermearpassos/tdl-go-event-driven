package command

import (
	"context"
	"tickets/domain"
)

type RecordTicketHandler struct {
	tracker domain.Tracker
}

func NewRecordTicketHandler(tracker domain.Tracker) RecordTicketHandler {
	return RecordTicketHandler{tracker: tracker}
}

func (h RecordTicketHandler) Handle(ctx context.Context, event domain.TicketBookingConfirmed) error {
	return h.tracker.AppendRow(ctx, "tickets-to-print", []string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency})
}

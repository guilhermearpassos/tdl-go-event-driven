package command

import (
	"context"
	"tickets/domain"
)

type RecordRefundHandler struct {
	tracker domain.Tracker
}

func NewRecordRefundHandler(tracker domain.Tracker) RecordRefundHandler {
	return RecordRefundHandler{tracker: tracker}
}

func (h RecordRefundHandler) Handle(ctx context.Context, event domain.TicketBookingCanceled) error {
	return h.tracker.AppendRow(ctx, "tickets-to-refund", []string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency})
}

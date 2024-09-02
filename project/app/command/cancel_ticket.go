package command

import (
	"context"
	"tickets/domain"
)

type CancelTicketHandler struct {
	repo domain.Repository
}

func NewCancelTicketHandler(repo domain.Repository) CancelTicketHandler {
	return CancelTicketHandler{repo: repo}
}

func (h *CancelTicketHandler) Handle(ctx context.Context, ticketID string) error {

	return h.repo.DeleteTicket(ctx, ticketID)
}

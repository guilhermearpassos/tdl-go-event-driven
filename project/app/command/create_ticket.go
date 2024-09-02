package command

import (
	"context"
	"tickets/domain"
)

type CreateTicket struct {
	TicketStatus domain.TicketStatus
}

type CreateTicketHandler struct {
	repo domain.Repository
}

func NewCreateTicketHandler(repo domain.Repository) CreateTicketHandler {
	return CreateTicketHandler{repo: repo}
}

func (h *CreateTicketHandler) Handle(ctx context.Context, cmd CreateTicket) error {

	return h.repo.CreateTicket(ctx, cmd.TicketStatus)
}

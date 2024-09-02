package query

import (
	"context"
	"tickets/domain"
)

type ListTicketsHandler struct {
	repo domain.Repository
}

func NewListTicketsHandler(repo domain.Repository) ListTicketsHandler {
	return ListTicketsHandler{repo: repo}
}

func (h ListTicketsHandler) Handle(ctx context.Context) ([]domain.TicketStatus, error) {
	return h.repo.GetTickets(ctx)
}

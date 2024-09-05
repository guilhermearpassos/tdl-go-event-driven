package command

import (
	"context"
	"tickets/domain"
)

type BookTickets struct {
	ShowId          string
	BookingId       string
	CustomerEmail   string
	NumberOfTickets int
}

type BookTicketsHandler struct {
	repo domain.Repository
}

func NewBookTicketsHandler(repo domain.Repository) BookTicketsHandler {
	return BookTicketsHandler{repo: repo}
}

func (h BookTicketsHandler) Handle(ctx context.Context, cmd BookTickets) error {
	return h.repo.CreateBooking(ctx, cmd.ShowId, cmd.BookingId, cmd.CustomerEmail, cmd.NumberOfTickets)
}

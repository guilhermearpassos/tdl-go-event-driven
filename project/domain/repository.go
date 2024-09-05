package domain

import "context"

type Repository interface {
	CreateTicket(ctx context.Context, ticket TicketStatus) error
	DeleteTicket(ctx context.Context, ticketID string) error
	GetTickets(ctx context.Context) ([]TicketStatus, error)
	CreateShow(ctx context.Context, show Show) error
	CreateBooking(ctx context.Context, showId string, bookingId string, customerEmail string, numberOfTickets int) error
}

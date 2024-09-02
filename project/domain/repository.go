package domain

import "context"

type Repository interface {
	CreateTicket(ctx context.Context, ticket TicketStatus) error
	DeleteTicket(ctx context.Context, ticketID string) error
	GetTickets(ctx context.Context) ([]TicketStatus, error)
}

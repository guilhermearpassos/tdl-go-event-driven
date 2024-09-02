package adapters

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"tickets/domain"
)

type PGTicketRepository struct {
	db *sqlx.DB
}

func NewPGTicketRepository(db *sqlx.DB) PGTicketRepository {
	return PGTicketRepository{db: db}
}

var _ domain.Repository = (*PGTicketRepository)(nil)

func (p PGTicketRepository) CreateTicket(ctx context.Context, ticket domain.TicketStatus) error {
	_, err := p.db.ExecContext(ctx, `
insert into tickets (ticket_id, price_amount, price_currency, customer_email) values ($1, $2, $3, $4) ON CONFLICT DO NOTHING

`, ticket.TicketID, ticket.Price.Amount, ticket.Price.Currency, ticket.CustomerEmail)
	if err != nil {
		return fmt.Errorf("creating ticket %s: %w", ticket.TicketID, err)
	}
	return nil
}

func (p PGTicketRepository) DeleteTicket(ctx context.Context, ticketID string) error {
	_, err := p.db.ExecContext(ctx, `DELETE FROM tickets where ticket_id=$1`, ticketID)
	if err != nil {
		return fmt.Errorf("deleting ticket %s: %w", ticketID, err)
	}
	return nil

}

func (p PGTicketRepository) GetTickets(ctx context.Context) ([]domain.TicketStatus, error) {
	rows, err := p.db.QueryxContext(ctx, `select ticket_id, price_amount, price_currency, customer_email from tickets`)
	if err != nil {
		return nil, fmt.Errorf("getting tickets queryCtx: %w", err)
	}
	tickets := make([]domain.TicketStatus, 0)
	for rows.Next() {
		var ticketID string
		var priceAmount float64
		var priceCurrency string
		var customerEmail string
		if err := rows.Scan(&ticketID, &priceAmount, &priceCurrency, &customerEmail); err != nil {
			return nil, fmt.Errorf("scanning ticket row: %w", err)
		}
		tickets = append(tickets, domain.TicketStatus{
			TicketID: ticketID,
			Status:   "confirmed",
			Price: domain.Money{
				Amount:   fmt.Sprintf("%.2f", priceAmount),
				Currency: priceCurrency,
			},
			CustomerEmail: customerEmail,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ticket rows err: %w", err)
	}
	return tickets, nil
}

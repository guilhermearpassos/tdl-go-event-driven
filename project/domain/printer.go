package domain

import "context"

type Printer interface {
	PrintTicket(ctx context.Context, name string, text string) error
}

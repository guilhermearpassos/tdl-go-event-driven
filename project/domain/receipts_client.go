package domain

import (
	"context"
)

type IssueReceiptRequest struct {
	TicketID       string
	Price          Money
	IdempotencyKey string
}

type ReceiptIssuer interface {
	IssueReceipt(ctx context.Context, receipt IssueReceiptRequest) error
}

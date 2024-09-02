package domain

import (
	"context"
)

type IssueReceiptRequest struct {
	TicketID string
	Price    Money
}

type ReceiptIssuer interface {
	IssueReceipt(ctx context.Context, receipt IssueReceiptRequest) error
}

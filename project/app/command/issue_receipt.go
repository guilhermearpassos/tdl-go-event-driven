package command

import (
	"context"
	"tickets/domain"
)

type IssueReceiptHandler struct {
	receiptClient domain.ReceiptIssuer
}

func NewIssueReceiptHandler(receiptClient domain.ReceiptIssuer) IssueReceiptHandler {
	return IssueReceiptHandler{receiptClient: receiptClient}
}

func (h IssueReceiptHandler) Handle(ctx context.Context, cmd domain.IssueReceiptRequest) error {
	return h.receiptClient.IssueReceipt(ctx, cmd)
}

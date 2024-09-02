package mock_services

import (
	"context"
	"sync"
	"tickets/domain"
)

type ReceiptsServiceMock struct {
	IssuedReceipts map[string]domain.IssueReceiptRequest
	mutex          *sync.Mutex
}

func NewReceiptsServiceMock() ReceiptsServiceMock {
	return ReceiptsServiceMock{
		IssuedReceipts: make(map[string]domain.IssueReceiptRequest),
		mutex:          &sync.Mutex{},
	}
}

var _ domain.ReceiptIssuer = (*ReceiptsServiceMock)(nil)

func (r *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request domain.IssueReceiptRequest) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.IssuedReceipts[request.TicketID] = request
	return nil
}

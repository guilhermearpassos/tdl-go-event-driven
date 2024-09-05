package mock_services

import (
	"context"
	"tickets/domain"
)

type MockPrinterClient struct {
	printedFiles map[string]string
}

func NewMockPrinterClient() MockPrinterClient {
	return MockPrinterClient{printedFiles: make(map[string]string)}
}

var _ domain.Printer = (*MockPrinterClient)(nil)

func (m MockPrinterClient) PrintTicket(ctx context.Context, name string, text string) error {
	if _, ok := m.printedFiles[name]; !ok {
		m.printedFiles[name] = text
	}
	return nil

}

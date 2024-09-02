package adapters

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"strconv"
	"tickets/domain"
)

type PrinterClientAdapter struct {
	clients *clients.Clients
}

func NewPrinterClientAdapter(clients *clients.Clients) *PrinterClientAdapter {
	return &PrinterClientAdapter{clients: clients}
}

var _ domain.Printer = (*PrinterClientAdapter)(nil)

func (p PrinterClientAdapter) PrintTicket(ctx context.Context, name string, text string) error {
	resp, err := p.clients.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, name, text)
	if err != nil {
		return fmt.Errorf("saving file %s: %w", name, err)
	}
	if resp.StatusCode() == 409 {
		return nil
	}
	if strconv.Itoa(resp.StatusCode())[0] != '2' {
		return fmt.Errorf("saving file %s: %w", name, fmt.Errorf("unexpected status code: %d", resp.StatusCode()))
	}
	return nil
}

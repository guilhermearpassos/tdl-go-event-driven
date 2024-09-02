package command

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"tickets/domain"
	"time"
)

type PrintTicket struct {
	Ticket domain.TicketStatus
}

type PrintTicketHandler struct {
	printer  domain.Printer
	template string
	eventBus *cqrs.EventBus
}

func NewPrintTicketHandler(printer domain.Printer, eventBus *cqrs.EventBus) PrintTicketHandler {
	// language=HTML
	template := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Ticket %s details</title>
</head>
<body>
<div>TicketID</div><div>%s</div>
<br>
<div>Customer</div><div>%s</div>
<br>
<div>Price</div><div>%s %s</div>
<br>
</body>
</html>`
	return PrintTicketHandler{printer: printer, template: template, eventBus: eventBus}
}

func (h PrintTicketHandler) Handle(ctx context.Context, cmd PrintTicket) error {
	fileName := fmt.Sprintf("%s-ticket.html", cmd.Ticket.TicketID)
	err := h.printer.PrintTicket(ctx, fileName,
		fmt.Sprintf(h.template, cmd.Ticket.TicketID, cmd.Ticket.TicketID, cmd.Ticket.CustomerEmail,
			cmd.Ticket.Price.Currency, cmd.Ticket.Price.Amount),
	)
	if err != nil {
		return fmt.Errorf("printing ticket %s: %w", cmd.Ticket.TicketID, err)
	}
	err = h.eventBus.Publish(ctx, domain.TicketPrinted{
		Header: domain.EventHeader{
			ID:          uuid.NewString(),
			PublishedAt: time.Now().UTC(),
		},
		TicketID: cmd.Ticket.TicketID,
		FileName: fileName,
	})
	if err != nil {
		return fmt.Errorf("publishing printed event %s: %w", cmd.Ticket.TicketID, err)
	}
	return nil
}

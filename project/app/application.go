package app

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"tickets/app/command"
	"tickets/app/query"
	"tickets/domain"
)

type Application struct {
	Queries  Queries
	Commands Commands
}
type Queries struct {
	ListTickets query.ListTicketsHandler
}
type Commands struct {
	IssueReceipt command.IssueReceiptHandler
	RecordTicket command.RecordTicketHandler
	RecordRefund command.RecordRefundHandler
	CreateTicket command.CreateTicketHandler
	CancelTicket command.CancelTicketHandler
	PrintTicket  command.PrintTicketHandler
	CreateShow   command.CreateShowHandler
	BookTickets  command.BookTicketsHandler
}

func NewApplication(receiptsClient domain.ReceiptIssuer, spreadsheetsClient domain.Tracker,
	repo domain.Repository, printer domain.Printer, eventBus *cqrs.EventBus) Application {
	application := Application{
		Queries: Queries{
			ListTickets: query.NewListTicketsHandler(repo),
		},
		Commands: Commands{
			IssueReceipt: command.NewIssueReceiptHandler(receiptsClient),
			RecordTicket: command.NewRecordTicketHandler(spreadsheetsClient),
			RecordRefund: command.NewRecordRefundHandler(spreadsheetsClient),
			CreateTicket: command.NewCreateTicketHandler(repo),
			CancelTicket: command.NewCancelTicketHandler(repo),
			PrintTicket:  command.NewPrintTicketHandler(printer, eventBus),
			CreateShow:   command.NewCreateShowHandler(repo),
			BookTickets:  command.NewBookTicketsHandler(repo),
		},
	}
	return application
}

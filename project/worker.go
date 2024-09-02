package main

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

//
//type Message struct {
//	Task   Task
//	Ticket domain.Ticket
//}
//type Worker struct {
//	queue              chan Message
//	receiptsClient     domain.ReceiptsClient
//	spreadsheetsClient domain.SpreadsheetsClient
//}
//
//func NewWorker(
//	receiptsClient domain.ReceiptsClient,
//	spreadsheetsClient domain.SpreadsheetsClient,
//) *Worker {
//	return &Worker{
//		queue:              make(chan Message, 100),
//		receiptsClient:     receiptsClient,
//		spreadsheetsClient: spreadsheetsClient,
//	}
//}
//func (w *Worker) Send(msg ...Message) {
//	for _, m := range msg {
//		w.queue <- m
//	}
//}
//func (w *Worker) Run() {
//	for msg := range w.queue {
//		switch msg.Task {
//		case TaskIssueReceipt:
//			err := w.receiptsClient.IssueReceipt(context.Background(), msg.Ticket)
//			if err != nil {
//				logrus.Errorf("failed to issue receipt for ticketId %s : %v", msg.Ticket.TicketId, err)
//				w.Send(msg)
//			}
//		case TaskAppendToTracker:
//			err := w.spreadsheetsClient.AppendRow(context.Background(), "tickets-to-print", []string{msg.Ticket.TicketId, msg.Ticket.CustomerEmail, msg.Ticket.Price.Amount, msg.Ticket.Price.Currency})
//			if err != nil {
//				logrus.Errorf("failed to append ticket to tracker, ticketId %s : %v", msg.Ticket.TicketId, err)
//				w.Send(msg)
//			}
//		}
//	}
//}

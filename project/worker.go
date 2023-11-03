package main

import (
	"context"
	"github.com/sirupsen/logrus"
	main2 "tickets/domain"
)

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

type Message struct {
	Task     Task
	TicketID string
}
type Worker struct {
	queue              chan Message
	receiptsClient     main2.ReceiptsClient
	spreadsheetsClient main2.SpreadsheetsClient
}

func NewWorker(
	receiptsClient main2.ReceiptsClient,
	spreadsheetsClient main2.SpreadsheetsClient,
) *Worker {
	return &Worker{
		queue:              make(chan Message, 100),
		receiptsClient:     receiptsClient,
		spreadsheetsClient: spreadsheetsClient,
	}
}
func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		w.queue <- m
	}
}
func (w *Worker) Run() {
	for msg := range w.queue {
		switch msg.Task {
		case TaskIssueReceipt:
			err := w.receiptsClient.IssueReceipt(context.Background(), msg.TicketID)
			if err != nil {
				logrus.Errorf("failed to issue receipt for ticketId %s : %v", msg.TicketID, err)
				w.Send(msg)
			}
		case TaskAppendToTracker:
			err := w.spreadsheetsClient.AppendRow(context.Background(), "tickets-to-print", []string{msg.TicketID})
			if err != nil {
				logrus.Errorf("failed to append ticket to tracker, ticketId %s : %v", msg.TicketID, err)
				w.Send(msg)
			}
		}
	}
}

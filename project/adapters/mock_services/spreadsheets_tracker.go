package mock_services

import (
	"context"
)

type SpreadsheetsTrackerClient struct {
	RowsCreatedBySheet map[string][][]string
}

func NewSpreadsheetsClient() SpreadsheetsTrackerClient {
	return SpreadsheetsTrackerClient{
		RowsCreatedBySheet: make(map[string][][]string),
	}
}

func (c SpreadsheetsTrackerClient) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	if _, ok := c.RowsCreatedBySheet[spreadsheetName]; !ok {
		c.RowsCreatedBySheet[spreadsheetName] = make([][]string, 0)
	}
	c.RowsCreatedBySheet[spreadsheetName] = append(c.RowsCreatedBySheet[spreadsheetName], row)

	return nil
}

package domain

import (
	"context"
)

type Tracker interface {
	AppendRow(ctx context.Context, bucketName string, row []string) error
}

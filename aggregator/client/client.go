package client

import (
	"context"

	"github.com/mhg14/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}

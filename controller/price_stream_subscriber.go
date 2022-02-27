package controller

import (
	"context"
)

type PriceStreamSubscriber interface {
	Start(ctx context.Context)
	Stop()
	Destroy()
}

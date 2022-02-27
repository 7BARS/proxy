package controller

import (
	"context"
)

type MockPriceStreaming struct {
}

func NewMockPriceStreaming() *MockPriceStreaming {
	return &MockPriceStreaming{}
}

func (s *MockPriceStreaming) Start(ctx context.Context) {}
func (s *MockPriceStreaming) Stop()                     {}
func (s *MockPriceStreaming) Destroy()                  {}

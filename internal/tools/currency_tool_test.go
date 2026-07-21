package tools

import (
	"context"
	"testing"
	"time"

	"google.golang.org/adk/v2/agent"
)

type mockAgentContext struct {
	agent.Context
	stdCtx context.Context
}

func (m mockAgentContext) Deadline() (time.Time, bool) { return m.stdCtx.Deadline() }
func (m mockAgentContext) Done() <-chan struct{}       { return m.stdCtx.Done() }
func (m mockAgentContext) Err() error                  { return m.stdCtx.Err() }
func (m mockAgentContext) Value(key any) any           { return m.stdCtx.Value(key) }

func TestGetCurrencyTool(t *testing.T) {
	cTool, err := GetCurrencyTool()
	if err != nil {
		t.Fatalf("failed to create currency tool: %v", err)
	}

	if cTool.Name() != "convert_currency" {
		t.Errorf("expected tool name 'convert_currency', got %q", cTool.Name())
	}

	if cTool.Description() == "" {
		t.Errorf("expected non-empty description for currency tool")
	}
}

func TestConvertCurrency(t *testing.T) {
	mockCtx := mockAgentContext{stdCtx: context.Background()}

	tests := []struct {
		name         string
		args         *ConvertCurrencyArgs
		expectErr    bool
		expectedAmt  float64
	}{
		{
			name:        "USD to EUR",
			args:        &ConvertCurrencyArgs{Amount: 100, From: "USD", To: "EUR"},
			expectErr:   false,
			expectedAmt: 92.0,
		},
		{
			name:        "EUR to USD",
			args:        &ConvertCurrencyArgs{Amount: 92, From: "EUR", To: "USD"},
			expectErr:   false,
			expectedAmt: 100.0,
		},
		{
			name:        "Same currency (USD to USD)",
			args:        &ConvertCurrencyArgs{Amount: 50, From: "USD", To: "USD"},
			expectErr:   false,
			expectedAmt: 50.0,
		},
		{
			name:      "Nil args",
			args:      nil,
			expectErr: true,
		},
		{
			name:      "Zero amount",
			args:      &ConvertCurrencyArgs{Amount: 0, From: "USD", To: "EUR"},
			expectErr: true,
		},
		{
			name:      "Unsupported source currency",
			args:      &ConvertCurrencyArgs{Amount: 100, From: "XYZ", To: "USD"},
			expectErr: true,
		},
		{
			name:      "Unsupported target currency",
			args:      &ConvertCurrencyArgs{Amount: 100, From: "USD", To: "XYZ"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ConvertCurrency(mockCtx, tt.args)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.ConvertedAmount != tt.expectedAmt {
				t.Errorf("expected converted amount %f, got %f", tt.expectedAmt, res.ConvertedAmount)
			}
		})
	}
}

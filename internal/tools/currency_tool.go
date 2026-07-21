package tools

import (
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

var currencyTracer = otel.Tracer("currency-tool")

type ConvertCurrencyArgs struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

type ConvertCurrencyResults struct {
	OriginalAmount  float64 `json:"original_amount"`
	FromCurrency    string  `json:"from_currency"`
	ConvertedAmount float64 `json:"converted_amount"`
	ToCurrency      string  `json:"to_currency"`
	Rate            float64 `json:"rate"`
}

// Standard exchange rates relative to 1 USD for demonstration/exploration.
var usdExchangeRates = map[string]float64{
	"USD": 1.0,
	"EUR": 0.92,
	"GBP": 0.79,
	"JPY": 155.5,
	"CAD": 1.36,
	"AUD": 1.52,
	"INR": 83.4,
}

func ConvertCurrency(ctx agent.Context, args *ConvertCurrencyArgs) (*ConvertCurrencyResults, error) {
	_, span := currencyTracer.Start(ctx, "ConvertCurrency")
	defer span.End()

	if args == nil {
		err := errors.New("arguments are required")
		span.RecordError(err)
		return nil, err
	}

	from := strings.ToUpper(strings.TrimSpace(args.From))
	to := strings.ToUpper(strings.TrimSpace(args.To))

	span.SetAttributes(
		attribute.Float64("currency.amount", args.Amount),
		attribute.String("currency.from", from),
		attribute.String("currency.to", to),
	)

	if args.Amount <= 0 {
		err := errors.New("amount must be greater than 0")
		span.RecordError(err)
		return nil, err
	}

	fromRate, fromOk := usdExchangeRates[from]
	if !fromOk {
		err := fmt.Errorf("unsupported source currency: %s", from)
		span.RecordError(err)
		return nil, err
	}

	toRate, toOk := usdExchangeRates[to]
	if !toOk {
		err := fmt.Errorf("unsupported target currency: %s", to)
		span.RecordError(err)
		return nil, err
	}

	rate := toRate / fromRate
	convertedAmount := args.Amount * rate

	return &ConvertCurrencyResults{
		OriginalAmount:  args.Amount,
		FromCurrency:    from,
		ConvertedAmount: convertedAmount,
		ToCurrency:      to,
		Rate:            rate,
	}, nil
}

func GetCurrencyTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name: "convert_currency",
			Description: `Converts an amount from one currency to another (e.g. USD, EUR, GBP, JPY, CAD, AUD, INR).
Args:
amount (number): The amount of money to convert.
from (string): The source 3-letter currency code (e.g. USD, EUR, GBP, JPY).
to (string): The target 3-letter currency code (e.g. USD, EUR, GBP, JPY).
Returns:
The converted amount, target currency, and exchange rate applied.`,
		},
		ConvertCurrency,
	)
}

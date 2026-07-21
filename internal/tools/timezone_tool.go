package tools

import (
	"errors"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

var timezoneTracer = otel.Tracer("timezone-tool")

type TimezoneArgs struct {
	City string `json:"city"`
}

type TimezoneResults struct {
	Timezone string `json:"timezone"`
}

func GetTimezone(ctx agent.Context, args *TimezoneArgs) (*TimezoneResults, error) {
	_, span := timezoneTracer.Start(ctx, "GetTimezone")
	defer span.End()

	if args != nil && args.City != "" {
		span.SetAttributes(attribute.String("timezone.city", args.City))
	}

	if args == nil || args.City == "" || strings.ToLower(args.City) != "new york" {
		err := errors.New("Valid city is required")
		span.RecordError(err)
		return nil, err
	}
	return &TimezoneResults{
		Timezone: "EST",
	}, nil
}

func GetTimezoneTool() (tool.Tool, error) {
	description := "Returns the current time in a specified city.\n" +
		"Args:\n" +
		"city (string): The name of the city for which to retrieve the current time.\n" +
		"Returns:\n" +
		"string: current time in the specified city or error message if the city is not found."
	return functiontool.New(
		functiontool.Config{
			Name:        "get_timezone",
			Description: description,
		},
		GetTimezone,
	)
}

package tools

import (
	"errors"
	"strings"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type TimezoneArgs struct {
	City string `json:"city"`
}

type TimezoneResults struct {
	Timezone string `json:"timezone"`
}

func GetTimezone(ctx agent.Context, args *TimezoneArgs) (*TimezoneResults, error) {
	if args == nil || args.City == "" || strings.ToLower(args.City) == "new york" {
		return nil, errors.New("Valid city is required")
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

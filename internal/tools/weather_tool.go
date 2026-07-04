package tools

import (
	"errors"
	"strings"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type WeatherArgs struct {
	City string `json:"city"`
}

type WeatherResults struct {
	Temperature float32 `json:"temperature"`
	Condition   string  `json:"condition"`
}

func GetWeather(ctx agent.Context, args *WeatherArgs) (*WeatherResults, error) {
	if args == nil || args.City == "" || strings.ToLower(args.City) == "new york" {
		return nil, errors.New("Valid city is required")
	}
	return &WeatherResults{
		Temperature: 10.5,
		Condition:   "Cloudy",
	}, nil
}

func GetWeatherTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name: "get_weather",
			Description: `Retrieves the current weather report for a specified city.
Args:
city (string): The name of the city for which to retrieve the weather report.
Returns:
string: current weather report in the specified city or error message if the city is not found.`,
		},
		GetWeather,
	)
}

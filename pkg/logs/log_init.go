package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func NewLog(logLevel string) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	time_zone := os.Getenv("APP_TIMEZONE")
	location, err := time.LoadLocation(time_zone)
	if err != nil {
		location = time.UTC // fallback
	}

	// Human-readable console writer
	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stderr,
		FormatTimestamp: func(i interface{}) string {
			if t, ok := i.(string); ok {
				parsed, err := time.Parse(time.RFC3339, t)
				if err == nil {
					return parsed.In(location).Format("2006-01-02T15:04:05-07:00")
				}
			}
			return fmt.Sprintf("%v", i)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("\n%s", i)
		},
		FormatFieldName: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("\n%-13s:", i))
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
	}

	// JSON file writer
	jsonFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not open JSON log file: %v", err))
	}

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, jsonFile)

	Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
}

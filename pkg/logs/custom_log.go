package logs

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog"
)

// CustomLog represents a custom log entry with additional context.
type CustomLog struct {
	MessageID string
	LogReason string
	Function  string
	File      string
	Line      int
}

// LogToString returns a formatted log message string.
func (e *CustomLog) LogToString() string {
	return fmt.Sprintf("MessageID: %s, LogReason: %s, Function: %s, File: %s, Line: %d",
		e.MessageID, e.LogReason, e.Function, e.File, e.Line)
}

// NewCustomLog creates a new CustomLog with caller information and logs it based on the specified level.
func NewCustomLog(messageID string, logDesc string, logType ...string) *CustomLog {
	pc, file, line, ok := runtime.Caller(1)
	function := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			function = fn.Name()
		}
	}

	msg := &CustomLog{
		MessageID: messageID,
		LogReason: logDesc,
		Function:  function,
		File:      file,
		Line:      line,
	}

	// Safely extract log type from variadic argument
	logLevel := "info"
	if len(logType) > 0 {
		logLevel = logType[0]
	}
	level, levelIcon := resolveLogLevel(logLevel)

	Logger.WithLevel(level).
		Timestamp().
		Str("🔖 MessageID", msg.MessageID).
		Str("📝 Message", msg.LogReason).
		Str("📁 File", msg.File).
		Str("🔧 Function", msg.Function).
		Int("🔢 Line", msg.Line).
		Msg(levelIcon)

	return msg
}

func resolveLogLevel(levelStr string) (zerolog.Level, string) {
	switch levelStr {
	case "fatal":
		return zerolog.FatalLevel, "☠️ FATAL"
	case "error":
		return zerolog.ErrorLevel, "🛑 ERROR"
	case "warn":
		return zerolog.WarnLevel, "⚠️ WARN"
	case "info":
		return zerolog.InfoLevel, "ℹ️ INFO"
	case "debug":
		return zerolog.DebugLevel, "🐞 DEBUG"
	default:
		return zerolog.InfoLevel, "ℹ️ INFO"
	}
}

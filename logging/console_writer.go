package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

const (
	cReset    = 0
	cBold     = 1
	cRed      = 31
	cGreen    = 32
	cYellow   = 33
	cBlue     = 34
	cMagenta  = 35
	cCyan     = 36
	cGray     = 37
	cDarkGray = 90
)

var consoleBufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 100))
	},
}

// ConsoleWriter reads a JSON object per write operation and output an
// optionally colored human readable version on the Out writer.
type ConsoleWriter struct {
	Out     io.Writer
	NoColor bool
}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}
	err = json.Unmarshal(p, &event)
	if err != nil {
		return
	}
	buf := consoleBufPool.Get().(*bytes.Buffer)
	defer consoleBufPool.Put(buf)
	getField := func(fieldName string) string {
		if field, ok := event[fieldName].(string); ok {
			return field
		}
		return ""
	}
	lvlColor := cReset
	level := getField(zerolog.LevelFieldName)
	if level != "" {
		if !w.NoColor {
			lvlColor = levelColor(level)
		}
		//level = strings.ToUpper(level)[0:4]
		level = strings.ToUpper(level)
	} else {
		level = "????"
	}
	from := getField(fromFieldName)
	threadName := getField(threadNameFieldName)
	loggerName := getField(loggerFieldName)
	fmt.Fprintf(buf, "%s [%s] [%s] %s [%s] - %s\n",
		colorize(event[zerolog.TimestampFieldName], cDarkGray, !w.NoColor),
		colorize(from, cCyan, !w.NoColor),
		colorize(threadName, cReset, !w.NoColor),
		colorize(level, lvlColor, !w.NoColor),
		colorize(loggerName, cReset, !w.NoColor),
		colorize(event[zerolog.MessageFieldName], cBold, !w.NoColor))
	//    buf.WriteByte('\n')
	buf.WriteTo(w.Out)
	n = len(p)
	return
}

func colorize(s interface{}, color int, enabled bool) string {
	if !enabled {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, s)
}

func levelColor(level string) int {
	switch level {
	case "debug":
		return cMagenta
	case "info":
		return cGreen
	case "warn":
		return cYellow
	case "error", "fatal", "panic":
		return cRed
	default:
		return cReset
	}
}

func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

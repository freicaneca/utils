package logging

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	loggerFieldName     = "logger_name"
	fromFieldName       = "from"
	threadNameFieldName = "thread_name"
	appFieldName        = "app"
)

type Logger struct {
	logger zerolog.Logger
}

func init() {
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.MessageFieldName = "message"
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000000"
}

// {"app":"PARAMS","severity":0,"level":"INFO","source":"params01","message":"Iniciando transacoes get_parametros_vendedor (Http)","priority":0,"client_id":"24520","tags":["redepos","params","offline","new","_grokparsefailure_sysloginput"],"@timestamp":"2017-12-19T01:39:58.047Z","HOSTNAME":"app01","thread_name":"New I/O worker #5","level_value":20000,"service":"params","short_ts":"22:39:58.046","logdate":"2017-12-18","@version":1,"host":"172.16.15.3","logger_name":"ParamsTransactionHttp","facility":0,"severity_label":"Emergency","facility_label":"kernel","timestamp":"2017-12-18 22:39:58.046"}
// {"app":"GTW","severity":0,"level":"DEBUG","source":"gtw01","message":"[id: 0x4139e446, /191.21.55.1:4283 => /189.36.19.132:9074] [OPENED]","priority":0,"client_id":"","tags":["redepos","gtw","offline","new","_grokparsefailure_sysloginput"],"@timestamp":"2017-12-19T01:39:57.960Z","HOSTNAME":"gw01","thread_name":"New I/O server boss #17","level_value":10000,"service":"gtw","short_ts":"22:39:57.960","logdate":"2017-12-18","@version":1,"host":"189.36.19.132","logger_name":"GtwServerHandler","facility":0,"severity_label":"Emergency","facility_label":"kernel","timestamp":"2017-12-18 22:39:57.960"}

/*
<syslogHost>log01</syslogHost>
<includeCallerInfo>false</includeCallerInfo>
<customFields>{"app":"GTW"}</customFields>
*/

func getFuncName(name ...string) string {
	if len(name) > 0 {
		return name[0]
	}
	fname := getCallerFuncName(3)
	lindex := strings.LastIndex(fname, ".")
	if lindex != -1 {
		fname = fname[lindex+1:]
	}
	return fname
}

// TODO: App field must be set via config. file.
func New(name ...string) *Logger {
	loggerName := getFuncName(name...)
	l := new(Logger)
	// time         src     id      thread               level logger_name          message
	// 00:00:00.265 [gtw02] [22545] [New I/O worker #15] DEBUG [GtwServerDecoder] - BYTES: 168
	// "logger_name":"main","app":"dummy","client_id":"12345","thread_name":"New I/O server boss #17"
	l.logger = log.Output(ConsoleWriter{Out: os.Stderr, NoColor: false})
	l.logger = l.logger.With().
		Timestamp().
		Str(loggerFieldName, loggerName).
		Str(appFieldName, "").
		Str(fromFieldName, "").
		Str(threadNameFieldName, "").
		Logger()
	return l
}

func (l *Logger) New(name ...string) *Logger {
	loggerName := getFuncName(name...)
	nl := new(Logger)
	nl.logger = l.logger.With().Str(loggerFieldName, loggerName).Logger()
	return nl
}

func (l *Logger) SetString(field string, value string) *Logger {
	l.logger = l.logger.With().Str(field, value).Logger()
	return l
}

func (l *Logger) SetInt(field string, value int) *Logger {
	l.logger = l.logger.With().Int(field, value).Logger()
	return l
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.Info().Msg(fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Error().Msg(fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprintf(format, v...))
}

func (l *Logger) SetFrom(from string) *Logger {
	return l.SetString(fromFieldName, from)
}

func (l *Logger) SetThread(name string) *Logger {
	return l.SetString(threadNameFieldName, name)
}

func (l *Logger) SetApp(name string) *Logger {
	return l.SetString(appFieldName, name)
}

/*
 * Gets caller's function name.
 *
 * The level determines how deep it can go. By calling GetCallerFuncName() you
 * get the caller's function name. By calling GetFuncName(2) you get
 * the previous caller's function name.. and so on.
 */
func getCallerFuncName(level ...int) string {
	callerLevel := 1
	if len(level) > 0 {
		callerLevel = level[0]
	}
	// get caller's function name
	pc, _, _, _ := runtime.Caller(callerLevel)
	return runtime.FuncForPC(pc).Name()
}

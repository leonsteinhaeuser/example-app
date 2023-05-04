package log

import (
	"io"

	"github.com/rs/zerolog"
)

type Zerolog struct {
	module  string
	zerolog zerolog.Logger
}

type ZerologOption func(*zerolog.Logger)

// NewZerlog initializes a console logger.
// Additional options can be passed to the logger by specifying an option function.
func NewZerlog(options ...ZerologOption) *Zerolog {
	zl := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "2006-01-02 15:04:05Z07:00"
	})).With().Timestamp().Stack().Logger()
	for _, option := range options {
		option(&zl)
	}
	return &Zerolog{
		zerolog: zl,
	}
}

// NewZerologWithWriter initializes a logger with a writer.
// Additional options can be passed to the logger by specifying an option function.
func NewZerologWithWriter(writer io.Writer, options ...ZerologOption) *Zerolog {
	zl := zerolog.New(writer).With().Timestamp().Stack().Logger()
	for _, option := range options {
		option(&zl)
	}
	return &Zerolog{
		zerolog: zl,
	}
}

// SetLevel sets the log level.
func (z *Zerolog) SetLevel(level int8) Logger {
	z.zerolog = z.zerolog.Level(zerolog.Level(level))
	return z
}

// SetModule sets the module name.
func (z *Zerolog) SetModule(module string) Logger {
	z.module = module
	return z
}

// Trace adds a trace log to the message.
func (z *Zerolog) Trace() Field {
	return &zerologField{
		module: &z.module,
		event:  z.zerolog.Trace(),
		fields: make(map[string]interface{}),
	}
}

// Debug adds a debug log to the message.
func (z *Zerolog) Debug() Field {
	return &zerologField{
		module: &z.module,
		event:  z.zerolog.Debug(),
		fields: make(map[string]interface{}),
	}
}

// Info adds a info log to the message.
func (z *Zerolog) Info() Field {
	return &zerologField{
		module: &z.module,
		event:  z.zerolog.Info(),
		fields: make(map[string]interface{}),
	}
}

// Warn adds a warn log to the message.
func (z *Zerolog) Warn() Field {
	return &zerologField{
		module: &z.module,
		event:  z.zerolog.Warn(),
		fields: make(map[string]interface{}),
	}
}

// Error adds a error log to the message.
func (z *Zerolog) Error(err error) Field {
	return &zerologField{
		module: &z.module,
		event:  z.zerolog.Error(),
		err:    err,
		fields: make(map[string]interface{}),
	}
}

// Panic adds a panic log to the message.
func (z *Zerolog) Panic(err error) Field {
	return &zerologField{
		module: &z.module,
		event:  z.zerolog.Panic(),
		err:    err,
		fields: make(map[string]interface{}),
	}
}

type zerologField struct {
	fields map[string]interface{}
	event  *zerolog.Event
	module *string
	err    error
}

// Field adds a field to the log message.
func (z *zerologField) Field(key string, value interface{}) Field {
	z.fields[key] = value
	return z
}

// unset cleares field references.
func (z *zerologField) unset() {
	z.fields = nil
	z.event = nil
}

func (z *zerologField) Log(message string) {
	defer z.unset()
	if z.err != nil {
		z.event.Err(z.err)
	}
	if z.module != nil && *z.module != "" {
		z.Field("module", *z.module)
	}
	z.event.Fields(z.fields).Msg(message)
}

func (z *zerologField) Logf(format string, args ...interface{}) {
	defer z.unset()
	if z.err != nil {
		z.event.Err(z.err)
	}
	if z.module != nil && *z.module != "" {
		z.Field("module", *z.module)
	}
	z.event.Fields(z.fields).Msgf(format, args...)
}

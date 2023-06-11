package log

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

// FxLogger is an Fx event logger that logs events to Zero.
type FxLogger struct {
	Logger *zerolog.Logger
}

func (l *FxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Debug().Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			l.Logger.Debug().Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.Logger.Debug().Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("callee", e.FunctionName).
				Str("callee", e.CallerName).
				Msg("OnStop hook failed")
		} else {
			l.Logger.Debug().Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).Str("type", e.TypeName).Msg("Supplied")
		} else {
			l.Logger.Debug().Str("type", e.TypeName).Msg("Supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Debug().Str("type", rtype).
				Str("constructor", e.ConstructorName).
				Msg("Provided")
		}
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("Error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do nothing. Will log on Invoked.

	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Str("stack", e.Trace).
				Str("function", e.FunctionName).Msg("Invoke failed")
		} else {
			l.Logger.Debug().Str("function", e.FunctionName).Msg("Invoked")
		}
	case *fxevent.Stopping:
		l.Logger.Info().Str("signal", strings.ToUpper(e.Signal.String())).Msg("Received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("Stop failed")
		}
	case *fxevent.RollingBack:
		l.Logger.Error().Err(e.StartErr).Msg("Start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("Rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("Start failed")
		} else {
			l.Logger.Debug().Msg("Started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("Custom logger initialization failed")
		} else {
			l.Logger.Debug().Str("function", e.ConstructorName).Msg("Initialized custom fxevent.Logger")
		}
	}
}

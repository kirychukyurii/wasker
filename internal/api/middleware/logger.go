package middleware

import (
	"github.com/kirychukyurii/wasker/internal/pkg/handler"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
	"time"
)

// Log middlewares constants.
const (
	logID           = "@id"
	logUserId       = "@user_id"
	logRemoteIP     = "@remote_ip"
	logURI          = "@uri"
	logHost         = "@host"
	logMethod       = "@method"
	logPath         = "@path"
	logRoute        = "@route"
	logProtocol     = "@protocol"
	logReferer      = "@referer"
	logUserAgent    = "@user_agent"
	logStatus       = "@status"
	logError        = "@error"
	logLatency      = "@latency"
	logLatencyHuman = "@latency_human"
	logBytesIn      = "@bytes_in"
	logBytesOut     = "@bytes_out"
	logHeaderPrefix = "@header:"
	logQueryPrefix  = "@query:"
	logFormPrefix   = "@form:"
	logCookiePrefix = "@cookie:"
)

// string to int base conversion.
const base = 10

type LoggerMiddleware struct {
	handler handler.HttpHandler
	logger  logger.Logger
}

func NewLoggerMiddleware(handler handler.HttpHandler, logger logger.Logger) LoggerMiddleware {
	return LoggerMiddleware{
		handler: handler,
		logger:  logger,
	}
}

func (a LoggerMiddleware) Setup() {
	a.logger.Log.Info().Msg("Setting up logger middleware")
	a.handler.Engine.Use(a.core())
}

func (a LoggerMiddleware) core() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			fields := map[string]string{
				"user_id":    logUserId,
				"request_id": logID,
				"remote_ip":  logRemoteIP,
				"uri":        logURI,
				"query":      logQueryPrefix,
				"host":       logHost,
				"method":     logMethod,
				"status":     logStatus,
				"latency":    logLatencyHuman,
				"error":      logError,
			}

			logFields, err := mapFields(ctx, next, fields)

			a.logger.Log.Info().
				Fields(logFields).
				Msg("Handle request")

			return
		}
	}
}

// mapFields maps fields based on tag name.
func mapFields(ec echo.Context, h echo.HandlerFunc, fm map[string]string) (map[string]interface{}, error) {
	logFields := map[string]interface{}{}
	start := time.Now()

	err := h(ec)
	if err != nil {
		ec.Error(err)
	}

	elapsed := time.Since(start)
	tags := mapTags(ec, elapsed)

	if err != nil {
		tags[logError] = err
	}

	for k, tag := range fm {
		if tag == "" {
			continue
		}

		if value, ok := tags[tag]; ok {
			logFields[k] = value
			continue
		}

		switch {
		case strings.HasPrefix(tag, logHeaderPrefix):
			key := tag[len(logHeaderPrefix):]
			logFields[k] = ec.Request().Header.Get(key)
		case strings.HasPrefix(tag, logQueryPrefix):
			key := tag[len(logQueryPrefix):]
			logFields[k] = ec.QueryParam(key)
		case strings.HasPrefix(tag, logFormPrefix):
			key := tag[len(logFormPrefix):]
			logFields[k] = ec.FormValue(key)
		case strings.HasPrefix(tag, logCookiePrefix):
			key := tag[len(logCookiePrefix):]
			cookie, err := ec.Cookie(key)
			if err == nil {
				logFields[k] = cookie.Value
			}
		}
	}

	return logFields, err
}

// mapTags maps the log tags with its related data. Populate previously the
// key/value avoids the cyclomatic complexity of the log middlewares to
// identify each tag and value.
func mapTags(ec echo.Context, latency time.Duration) map[string]interface{} {
	tags := map[string]interface{}{}

	req := ec.Request()
	res := ec.Response()

	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}

	tags[logID] = id
	tags[logUserId] = ec.Get("user_id")
	tags[logRemoteIP] = ec.RealIP()
	tags[logURI] = req.RequestURI
	tags[logHost] = req.Host
	tags[logMethod] = req.Method

	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	tags[logPath] = path
	tags[logRoute] = ec.Path()
	tags[logProtocol] = req.Proto
	tags[logReferer] = req.Referer()
	tags[logUserAgent] = req.UserAgent()
	tags[logStatus] = res.Status
	tags[logLatency] = strconv.FormatInt(int64(latency), base)
	tags[logLatencyHuman] = latency.String()

	cl := req.Header.Get(echo.HeaderContentLength)
	if cl == "" {
		cl = "0"
	}

	tags[logBytesIn] = cl
	tags[logBytesOut] = strconv.FormatInt(res.Size, base)

	return tags
}

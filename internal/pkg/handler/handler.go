package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type HttpHandler struct {
	Engine   *echo.Echo
	RouterV1 *echo.Group

	Validate
}

// BinderWithValidation to verify the request's struct for parameter validation
type BinderWithValidation struct{}

func New(logger logger.Logger, config config.Config) HttpHandler {
	// Error handlers
	echo.NotFoundHandler = func(ctx echo.Context) error {
		return Response{Code: http.StatusNotFound}.JSON(ctx)
	}

	echo.MethodNotAllowedHandler = func(ctx echo.Context) error {
		return Response{Code: http.StatusMethodNotAllowed}.JSON(ctx)
	}

	// new Echo engine
	engine := echo.New()
	engine.HidePort = true
	engine.HideBanner = true
	engine.Binder = &BinderWithValidation{}

	// set http handler
	httpHandler := HttpHandler{
		Engine:   engine,
		RouterV1: engine.Group("/api/v1"),
	}

	// custom the error handler
	httpHandler.Engine.HTTPErrorHandler = func(err error, ctx echo.Context) {
		var (
			code    = http.StatusInternalServerError
			message interface{}
		)

		he, ok := err.(*echo.HTTPError)
		if ok {
			code = he.Code
			message = he.Message

			if he.Internal != nil {
				message = fmt.Errorf("%v - %v", message, he.Internal)
			}
		}

		// Send response
		if !ctx.Response().Committed {
			// https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html
			if ctx.Request().Method == http.MethodHead {
				err = ctx.NoContent(he.Code)
			} else {
				err = Response{Code: code, Message: message}.JSON(ctx)
			}

			if err != nil {
				logger.Log.Error().Err(err).Msg("Received method")
			}
		}
	}

	// override the default validator
	httpHandler.Engine.Validator = func() echo.Validator {
		v := validator.New()

		if err := v.RegisterValidation("json", func(fl validator.FieldLevel) bool {
			var js json.RawMessage
			return json.Unmarshal([]byte(fl.Field().String()), &js) == nil
		}); err != nil {
			logger.Log.Error().Err(err).Str("tag", "json").Msg("Failed to register validation")
		}

		if err := v.RegisterValidation("in", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if containsString(strings.Split(fl.Param(), ";"), value) || value == "" {
				return true
			}

			return false
		}); err != nil {
			logger.Log.Error().Err(err).Str("tag", "in").Msg("Failed to register validation")
		}

		return &Validator{
			validate: v,
		}
	}()

	return httpHandler
}

func containsString(items []string, item string) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

func (BinderWithValidation) Bind(i interface{}, ctx echo.Context) error {
	binder := &echo.DefaultBinder{}
	if err := binder.Bind(i, ctx); err != nil {
		return errors.New(err.(*echo.HTTPError).Message.(string))
	}

	if err := ctx.Validate(i); err != nil {
		// Validate only provides verification function for struct.
		// When the requested data type is not struct,
		// the variable should be considered legal after the bind succeeds.
		//if reflect.TypeOf(i).Kind() != reflect.Struct {
		//	return nil
		//}

		var buf bytes.Buffer
		if ferrs, ok := err.(validator.ValidationErrors); ok {
			buf.WriteString("validation failed on")

			for i, ferr := range ferrs {
				buf.WriteString(" " + ferr.Field() + " ")
				buf.WriteString("(" + ferr.Tag() + ")")
				if i != len(ferrs)-1 {
					buf.WriteString(",")
				}
			}

			return errors.New(buf.String())
		}

		return err
	}

	return nil
}

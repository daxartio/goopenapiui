package echoopenapiui

import (
	"github.com/daxartio/goopenapiui"
	"github.com/labstack/echo/v4"
)

func New(openapiui *goopenapiui.Openapiui) echo.MiddlewareFunc {
	handle := openapiui.Handler()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			handle(ctx.Response(), ctx.Request())
			if ctx.Response().Committed {
				return nil
			}
			return next(ctx)
		}
	}
}

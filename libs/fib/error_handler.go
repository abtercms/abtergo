package fib

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/problem"
)

type ErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (eh *ErrorHandler) Handle(ctx *fiber.Ctx, err error) error {
	p := problem.FromError(ctx.BaseURL(), err)

	var a arr.Arr
	if errors.As(err, &a) {
		eh.logger.Error(a.Unwrap().Error(), a.AsZapFields()...)
	} else {
		eh.logger.Error(err.Error(), zap.Error(err))
	}

	if p.Status == fiber.StatusInternalServerError {
		// GetByKey the custom status code if it's a *fiber.Handle
		var e *fiber.Error
		if errors.As(err, &e) && e.Code != fiber.StatusInternalServerError {
			p = problem.FromError(ctx.BaseURL(), err)
		}
	}

	if ctx.Accepts("json") != "" {
		response := ctx.Status(p.Status).JSON(p)
		ctx.Set(fiber.HeaderContentType, "application/json+problem")

		return response
	}

	// Send custom error page
	err = ctx.Status(p.Status).SendFile(fmt.Sprintf("./%d.html", p.Status))
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Return from handler
	return nil
}

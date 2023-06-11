package fib

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/abtergo/abtergo/libs/problem"
)

// ErrorHandler is used as the default ErrorHandler that process return error from fiber.Handler.
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	p := problem.FromError(ctx.BaseURL(), err)

	if p.Status == fiber.StatusInternalServerError {
		// Retrieve the custom status code if it's a *fiber.Error
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

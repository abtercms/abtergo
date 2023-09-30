package website

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/abtergo/abtergo/libs/arr"
)

// Handler represents a set of HTTP handlers.
type Handler struct {
	logger  *slog.Logger
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// AddRoutes registers routes in a fiber.Router.
func (h *Handler) AddRoutes(r fiber.Router) {
	r.Use(h.CatchAll)
}

// CatchAll handles any requests which has not been caught by other handlers.
func (h *Handler) CatchAll(c *fiber.Ctx) error {
	h.logger.Info("catch all", slog.String("method", c.Method()), slog.String("path", c.Path()))

	if c.Method() != fiber.MethodGet {
		return c.SendStatus(fiber.StatusMethodNotAllowed)
	}

	statusCode := fiber.StatusOK
	body, err := h.service.Get(c.Context(), c.BaseURL(), c.Path())
	if err != nil {
		statusCode = arr.HTTPStatusFromError(err)
	}

	return c.Status(statusCode).SendString(body)
}

package website

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
)

// Handler represents a set of HTTP handlers.
type Handler struct {
	logger  *zap.Logger
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service, logger *zap.Logger) *Handler {
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
	h.logger.Info("catch all", zap.String("method", c.Method()), zap.String("path", c.Path()))

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

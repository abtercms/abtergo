package page

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/fib"
)

// Handler represents a set of HTTP handlers.
type Handler struct {
	logger  *zap.Logger
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(logger *zap.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// AddAPIRoutes registers routes in a fiber.Router.
func (h *Handler) AddAPIRoutes(r fiber.Router) {
	r.Post("/pages", h.Post)
	r.Get("/pages", h.List)
	r.Get("/pages/:id", h.Get)
	r.Put("/pages/:id", h.Put)
	r.Delete("/pages/:id", h.Delete)
	r.Post("/pages/:id/activations", h.Activate)
	r.Post("/pages/:id/inactivations", h.Inactivate)
}

// Post handles requests to persist new Page.
func (h *Handler) Post(c *fiber.Ctx) error {
	payload := Page{}

	if err := c.BodyParser(&payload); err != nil {
		return arr.Wrap(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	response, err := h.service.Create(c.Context(), payload)
	if err != nil {
		return errors.Wrap(err, "failed to create the page")
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// List handles requests to retrieve a collection of Page instances.
func (h *Handler) List(c *fiber.Ctx) error {
	filter := Filter{
		Website: c.Get("website"),
		Path:    c.Get("path"),
	}

	response, err := h.service.List(c.Context(), filter)
	if err != nil {
		return fmt.Errorf("failed to retrieve entity, err: %w", err)
	}

	return c.JSON(response)
}

// Get handles requests to retrieve a single Page instance.
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	response, err := h.service.Get(c.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to retrieve entity, err: %w", err)
	}

	return c.JSON(response)
}

// Put handles requests to update a single Page instance.
func (h *Handler) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	payload := Page{}

	if err := c.BodyParser(&payload); err != nil {
		return arr.Wrap(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	etag := c.Get(fiber.HeaderETag)

	response, err := h.service.Update(c.Context(), id, payload, etag)
	if err != nil {
		return errors.Wrap(err, "failed to update entity")
	}

	return c.JSON(response)
}

// Delete handles requests to delete a single Page instance.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	err := h.service.Delete(c.Context(), id, c.Get(fiber.HeaderETag))
	if err != nil {
		return fmt.Errorf("failed to delete entity, err: %w", err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Activate handles requests to activate a single Page instance.
func (h *Handler) Activate(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	etag := c.Get(fiber.HeaderETag)

	response, err := h.service.Transition(c.Context(), id, Activate, etag)
	if err != nil {
		return fmt.Errorf("failed to activate entity, err: %w", err)
	}

	return c.JSON(response)
}

// Inactivate handles requests to inactivate a single Page instance.
func (h *Handler) Inactivate(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	etag := c.Get(fiber.HeaderETag)

	response, err := h.service.Transition(c.Context(), id, Inactivate, etag)
	if err != nil {
		return fmt.Errorf("failed to inactivate entity, err: %w", err)
	}

	return c.JSON(response)
}

func checkWiring(id string) {
	if id == "" {
		panic(fib.ErrRouteHandleWiring)
	}
}

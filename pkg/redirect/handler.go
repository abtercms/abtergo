package redirect

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
	r.Post("/redirects", h.Post)
	r.Get("/redirects", h.List)
	r.Get("/redirects/:id", h.Get)
	r.Put("/redirects/:id", h.Put)
	r.Delete("/redirects/:id", h.Delete)
}

// Post handles requests to persist new Redirect.
func (h *Handler) Post(c *fiber.Ctx) error {
	payload := Redirect{}

	if err := c.BodyParser(&payload); err != nil {
		return arr.Wrap(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	response, err := h.service.Create(c.Context(), payload)
	if err != nil {
		return errors.Wrap(err, "failed to create entity")
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// List handles requests to retrieve a collection of Redirect instances.
func (h *Handler) List(c *fiber.Ctx) error {
	// TODO: Smarter filter parsing
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

// Get handles requests to retrieve a single Redirect instance.
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	response, err := h.service.Get(c.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to retrieve entity, err: %w", err)
	}

	return c.JSON(response)
}

// Put handles requests to update a single Redirect instance.
func (h *Handler) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	payload := Redirect{}

	if err := c.BodyParser(&payload); err != nil {
		return arr.Wrap(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	etag := c.Get(fiber.HeaderETag)

	response, err := h.service.Update(c.Context(), id, payload, etag)
	if err != nil {
		return fmt.Errorf("failed to update entity, err: %w", err)
	}

	return c.JSON(response)
}

// Delete handles requests to delete a single Handler instance.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	err := h.service.Delete(c.Context(), id, c.Get(fiber.HeaderETag))
	if err != nil {
		return fmt.Errorf("failed to delete entity, err: %w", err)
	}

	return c.SendStatus(http.StatusNoContent)
}

func checkWiring(id string) {
	if id == "" {
		panic(fib.ErrRouteHandleWiring)
	}
}

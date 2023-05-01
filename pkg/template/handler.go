package template

import (
	"fmt"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"

	"github.com/abtergo/abtergo/libs/ablog"
	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/pkg/common"
)

// Handler represents a set of HTTP handlers.
type Handler struct {
	logger  ablog.Logger
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(logger ablog.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// AddAPIRoutes registers routes in a fiber.Router.
func (h *Handler) AddAPIRoutes(r fiber.Router) {
	r.Post("/templates", h.Post)
	r.Get("/templates", h.List)
	r.Get("/templates/:id", h.Get)
	r.Put("/templates/:id", h.Put)
	r.Delete("/templates/:id", h.Delete)
}

// Post handles requests to persist new Template.
func (h *Handler) Post(c *fiber.Ctx) error {
	payload := Template{}

	if err := c.BodyParser(&payload); err != nil {
		return fmt.Errorf("failed to parse the request payload, err: %w", arr.Wrap(arr.InvalidUserInput, err))
	}

	response, err := h.service.Create(c.Context(), payload)
	if err != nil {
		return fmt.Errorf("failed to create entity, err: %w", err)
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// List handles requests to retrieve a collection of Template instances.
func (h *Handler) List(c *fiber.Ctx) error {
	// TODO: Smarter filter parsing
	filter := Filter{
		Website: c.Get("website"),
		Name:    c.Get("name"),
	}

	response, err := h.service.List(c.Context(), filter)
	if err != nil {
		return fmt.Errorf("failed to retrieve entity, err: %w", err)
	}

	return c.JSON(response)
}

// Get handles requests to retrieve a single Template instance.
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	response, err := h.service.Get(c.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to retrieve entity, err: %w", err)
	}

	return c.JSON(response)
}

// Put handles requests to update a single Template instance.
func (h *Handler) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	payload := Template{}

	if err := c.BodyParser(&payload); err != nil {
		return fmt.Errorf("failed to parse the request payload, err: %w", arr.Wrap(arr.InvalidUserInput, err))
	}

	etag := c.Get(fiber.HeaderETag)

	response, err := h.service.Update(c.Context(), id, payload, etag)
	if err != nil {
		return fmt.Errorf("failed to update entity, err: %w", err)
	}

	return c.JSON(response)
}

// Delete handles requests to delete a single Template instance.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	err := h.service.Delete(c.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to delete entity, err: %w", err)
	}

	return c.SendStatus(http.StatusNoContent)
}

func checkWiring(id string) {
	if id == "" {
		panic(common.ErrRouteHandleWiring)
	}
}

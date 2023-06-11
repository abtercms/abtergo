package block

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
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
	r.Post("/blocks", h.Post)
	r.Get("/blocks", h.List)
	r.Get("/blocks/:id", h.Get)
	r.Put("/blocks/:id", h.Put)
	r.Delete("/blocks/:id", h.Delete)
}

// Post handles requests to persist new Block.
func (h *Handler) Post(c *fiber.Ctx) error {
	payload := Block{}

	if err := c.BodyParser(&payload); err != nil {
		return arr.Wrap(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	response, err := h.service.Create(c.Context(), payload)
	if err != nil {
		return fmt.Errorf("failed to create entity, err: %w", err)
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// List handles requests to retrieve a collection of Block instances.
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

// Get handles requests to retrieve a single Block instance.
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	response, err := h.service.Get(c.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to retrieve entity, err: %w", err)
	}

	return c.JSON(response)
}

// Put handles requests to update a single Block instance.
func (h *Handler) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	payload := Block{}

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

// Delete handles requests to delete a single Block instance.
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

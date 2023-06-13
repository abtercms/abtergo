package template

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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
		return arr.WrapWithType(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	response, err := h.service.Create(c.Context(), payload)
	if err != nil {
		return errors.Wrap(err, "failed to create entity")
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// List handles requests to retrieve a collection of Template instances.
func (h *Handler) List(c *fiber.Ctx) error {
	filter := Filter{
		Website: c.Get("website"),
		Name:    c.Get("name"),
	}

	response, err := h.service.List(c.Context(), filter)
	if err != nil {
		return errors.Wrap(err, "failed to list entities")
	}

	return c.JSON(response)
}

// Get handles requests to retrieve a single Template instance.
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	response, err := h.service.Get(c.Context(), id)
	if err != nil {
		return errors.Wrap(err, "failed to get entity")
	}

	return c.JSON(response)
}

// Put handles requests to update a single Template instance.
func (h *Handler) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	payload := Template{}

	if err := c.BodyParser(&payload); err != nil {
		return arr.WrapWithType(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	etag := c.Get(fiber.HeaderETag)

	response, err := h.service.Update(c.Context(), id, payload, etag)
	if err != nil {
		return errors.Wrap(err, "failed to update entity")
	}

	return c.JSON(response)
}

// Delete handles requests to delete a single Template instance.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	err := h.service.Delete(c.Context(), id, c.Get(fiber.HeaderETag))
	if err != nil {
		return errors.Wrap(err, "failed to delete entity")
	}

	return c.SendStatus(http.StatusNoContent)
}

func checkWiring(id string) {
	if id == "" {
		panic(arr.New(arr.ApplicationError, "wiring error"))
	}
}

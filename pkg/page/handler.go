package page

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
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
		return arr.WrapWithType(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	if payload.ID != "" {
		return arr.New(arr.InvalidUserInput, "id provided", slog.String("id in payload", payload.ID.String()))
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
		return errors.Wrap(err, "failed to list the pages")
	}

	return c.JSON(response)
}

// Get handles requests to retrieve a single Page instance.
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	response, err := h.service.Get(c.Context(), model.ID(id))
	if err != nil {
		return errors.Wrap(err, "failed to get the page")
	}

	return c.JSON(response)
}

// Put handles requests to update a single Page instance.
func (h *Handler) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	eTag := c.Get(fiber.HeaderETag)
	checkWiring(id)

	if eTag == "" {
		return arr.New(arr.InvalidUserInput, "missing e-tag")
	}

	payload := Page{}

	if err := c.BodyParser(&payload); err != nil {
		return arr.WrapWithType(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	if payload.ID.String() != id {
		return arr.New(arr.InvalidUserInput, "id mismatch", slog.String("id in path", id), slog.String("id in payload", payload.ID.String()))
	}

	response, err := h.service.Update(c.Context(), payload, model.ETag(eTag))
	if err != nil {
		return errors.Wrap(err, "failed to update entity")
	}

	return c.JSON(response)
}

// Delete handles requests to delete a single Page instance.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	eTag := c.Get(fiber.HeaderETag)
	checkWiring(id)

	err := h.service.Delete(c.Context(), model.ID(id), model.ETag(eTag))
	if err != nil {
		return errors.Wrap(err, "failed to delete entity")
	}

	return c.SendStatus(http.StatusNoContent)
}

// Activate handles requests to activate a single Page instance.
func (h *Handler) Activate(c *fiber.Ctx) error {
	id := c.Params("id")
	eTag := c.Get(fiber.HeaderETag)
	checkWiring(id)

	response, err := h.service.Transition(c.Context(), model.ID(id), Activate, model.ETag(eTag))
	if err != nil {
		return errors.Wrap(err, "failed to activate entity")
	}

	return c.JSON(response)
}

// Inactivate handles requests to inactivate a single Page instance.
func (h *Handler) Inactivate(c *fiber.Ctx) error {
	id := c.Params("id")
	eTag := c.Get(fiber.HeaderETag)
	checkWiring(id)

	response, err := h.service.Transition(c.Context(), model.ID(id), Inactivate, model.ETag(eTag))
	if err != nil {
		return errors.Wrap(err, "failed to inactivate entity")
	}

	return c.JSON(response)
}

func checkWiring(id string) {
	if id == "" {
		panic(arr.New(arr.ApplicationError, "wiring error"))
	}
}

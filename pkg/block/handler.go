package block

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
		return arr.WrapWithType(arr.InvalidUserInput, err, "failed to parse the request payload")
	}

	if payload.ID != "" {
		return arr.New(arr.InvalidUserInput, "id provided", slog.String("id in payload", payload.ID.String()))
	}

	response, err := h.service.Create(c.Context(), payload)
	if err != nil {
		return errors.Wrap(err, "failed to create entity")
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// List handles requests to retrieve a collection of Block instances.
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

// Get handles requests to retrieve a single Block instance.
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	checkWiring(id)

	response, err := h.service.GetByID(c.Context(), model.ID(id))
	if err != nil {
		return errors.Wrap(err, "failed to get entity")
	}

	return c.JSON(response)
}

// Put handles requests to update a single Block instance.
func (h *Handler) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	eTag := c.Get(fiber.HeaderETag)
	checkWiring(id)

	if eTag == "" {
		return arr.New(arr.InvalidUserInput, "missing e-tag")
	}

	payload := Block{}

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

// Delete handles requests to delete a single Block instance.
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

func checkWiring(id string) {
	if id == "" {
		panic(arr.New(arr.ApplicationError, "wiring error"))
	}
}

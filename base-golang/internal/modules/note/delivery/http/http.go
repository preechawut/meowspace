package http

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"base-golang/internal/modules/note/domain"
	"base-golang/internal/modules/note/service"
	"base-golang/internal/platform/httpserver"
	"base-golang/internal/shared/apperror"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List(c fiber.Ctx) error {
	ctx := c.Context()

	options, err := decodeListOptions(c)
	if err != nil {
		return err
	}

	page, err := h.service.List(ctx, options)
	if err != nil {
		return err
	}

	return httpserver.JSON(c, fiber.StatusOK, NotePageResponse{
		Items:      toNoteResponses(page.Notes),
		Page:       options.Page,
		PageSize:   options.PageSize,
		Total:      page.Total,
		TotalPages: page.TotalPages,
	})
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input SaveNoteRequest
	err := httpserver.Decode(c, &input)
	if err != nil {
		return err
	}

	ctx := c.Context()
	result, err := h.service.Create(ctx, domain.Note{
		Title: input.Title,
	})
	if err != nil {
		return err
	}

	return httpserver.JSON(c, fiber.StatusCreated, toNoteResponse(result))
}

func decodeListOptions(c fiber.Ctx) (service.ListOptions, error) {
	options := service.ListOptions{
		Page:     1,
		PageSize: 20,
	}

	pageStr := c.Query("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return options, apperror.New(apperror.Invalid, "page must be an integer")
		}
		options.Page = page
	}

	pageSizeStr := c.Query("page_size")
	if pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return options, apperror.New(apperror.Invalid, "page_size must be an integer")
		}
		options.PageSize = pageSize
	}

	return options, nil
}

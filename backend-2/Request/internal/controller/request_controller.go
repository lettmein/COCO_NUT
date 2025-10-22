package controller

import (
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/dto"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/service"
	"github.com/gofiber/fiber/v2"
)

type RequestController struct {
	service *service.RequestService
}

func NewRequestController(service *service.RequestService) *RequestController {
	return &RequestController{service: service}
}

func (c *RequestController) CreateRequest(ctx *fiber.Ctx) error {
	var req dto.CreateRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	result, err := c.service.CreateRequest(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func (c *RequestController) GetRequest(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Request ID is required",
		})
	}

	result, err := c.service.GetRequest(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: "Request not found",
		})
	}

	return ctx.JSON(result)
}

func (c *RequestController) GetAllRequests(ctx *fiber.Ctx) error {
	status := ctx.Query("status")

	var results []*dto.RequestResponse
	var err error

	if status != "" {
		results, err = c.service.GetRequestsByStatus(status)
	} else {
		results, err = c.service.GetAllRequests()
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(results)
}

func (c *RequestController) UpdateRequestStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Request ID is required",
		})
	}

	var req dto.UpdateRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err := c.service.UpdateRequestStatus(id, req.Status)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(dto.SuccessResponse{
		Message: "Request status updated successfully",
	})
}

func (c *RequestController) DeleteRequest(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Request ID is required",
		})
	}

	err := c.service.DeleteRequest(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(dto.SuccessResponse{
		Message: "Request deleted successfully",
	})
}

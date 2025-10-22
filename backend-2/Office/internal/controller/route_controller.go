package controller

import (
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/dto"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/service"
	"github.com/gofiber/fiber/v2"
)

type RouteController struct {
	service *service.RouteService
}

func NewRouteController(service *service.RouteService) *RouteController {
	return &RouteController{service: service}
}

func (c *RouteController) CreateRoute(ctx *fiber.Ctx) error {
	var req dto.CreateRouteDTO
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	result, err := c.service.CreateRoute(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func (c *RouteController) GetRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Route ID is required",
		})
	}

	result, err := c.service.GetRoute(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: "Route not found",
		})
	}

	return ctx.JSON(result)
}

func (c *RouteController) GetAllRoutes(ctx *fiber.Ctx) error {
	status := ctx.Query("status")

	var results []*dto.RouteResponse
	var err error

	if status != "" {
		results, err = c.service.GetRoutesByStatus(status)
	} else {
		results, err = c.service.GetAllRoutes()
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(results)
}

func (c *RouteController) UpdateRouteStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Route ID is required",
		})
	}

	var req dto.UpdateRouteDTO
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err := c.service.UpdateRouteStatus(id, req.Status)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(dto.SuccessResponse{
		Message: "Route status updated successfully",
	})
}

func (c *RouteController) DeleteRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Route ID is required",
		})
	}

	err := c.service.DeleteRoute(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(dto.SuccessResponse{
		Message: "Route deleted successfully",
	})
}

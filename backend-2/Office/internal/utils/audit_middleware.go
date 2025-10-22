package utils

import (
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/client"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/dto"
	"github.com/gofiber/fiber/v2"
	"time"
)

func AuditMiddleware(serviceName string, auditClient *client.AuditClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestBody := string(c.Body())

		err := c.Next()

		processingTime := time.Since(start).Milliseconds()

		responseBody := string(c.Response().Body())

		auditLog := &dto.AuditLogDTO{
			ServiceName:    serviceName,
			RequestMethod:  c.Method(),
			RequestURL:     c.OriginalURL(),
			RequestBody:    requestBody,
			ResponseBody:   responseBody,
			StatusCode:     c.Response().StatusCode(),
			ProcessingTime: processingTime,
		}

		go auditClient.SendLog(auditLog)

		return err
	}
}

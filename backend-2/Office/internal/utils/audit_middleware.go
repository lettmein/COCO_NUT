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

		duration := time.Since(start)
		durationMs := float64(duration.Milliseconds())

		responseBody := string(c.Response().Body())

		auditLog := &dto.AuditLogDTO{
			ReqServiceType:  serviceName,
			RespServiceType: "audit-logger",
			Uri:             c.OriginalURL(),
			CreatedAt:       time.Now(),
			DurationTime:    durationMs,
			RequestBody:     requestBody,
			ResponseBody:    responseBody,
		}

		go auditClient.SendLog(auditLog)

		return err
	}
}

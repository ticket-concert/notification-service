package middleware

import (
	config "notification-service/configs"
	"notification-service/internal/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

type Middlewares struct {
	redisClient redis.Collections
}

func NewMiddlewares(redis redis.Collections) Middlewares {
	return Middlewares{
		redisClient: redis,
	}
}

func (m Middlewares) VerifyBasicAuth() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.GetConfig().UsernameBasicAuth: config.GetConfig().PasswordBasicAuth,
		},
	})
}

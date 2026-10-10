package app

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3"

	"base-golang/internal/config"
	notedelivery "base-golang/internal/modules/note/delivery/http"
	"base-golang/internal/platform/httpserver"
	"base-golang/internal/platform/httpserver/middleware"
)

func newRouter(db *sql.DB, cfg config.Config, note *notedelivery.Handler) *fiber.App {
	router := httpserver.New(cfg.TrustedProxies)

	// middleware
	router.Use(httpserver.ResponseHeaders())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logging())
	router.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	router.Use(middleware.RateLimit(cfg.APIRateLimit, time.Minute))
	router.Use(httpserver.RequestTimeout(15 * time.Second))
	router.Use(middleware.Recovery())

	// health checks / probe
	router.Get("/healthz", func(c fiber.Ctx) error {
		return httpserver.JSON(c, fiber.StatusOK, fiber.Map{"status": "ok"})
	})

	router.Get("/readyz", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			return httpserver.ErrorStatus(c, fiber.StatusServiceUnavailable, "database_unavailable", "database is unavailable")
		}

		return httpserver.JSON(c, fiber.StatusOK, fiber.Map{"status": "ready"})
	})

	// api v1
	v1 := router.Group("/api/v1")
	v1.Get("/notes", note.List)
	v1.Post("/notes", note.Create)

	return router
}

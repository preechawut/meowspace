package app

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"base-golang/internal/config"
	notedelivery "base-golang/internal/modules/note/delivery/http"
	noterepo "base-golang/internal/modules/note/infrastructure"
	noteservice "base-golang/internal/modules/note/service"
)

func New(db *sql.DB, cfg config.Config) *fiber.App {
	noteRepository := noterepo.New(db)
	noteService := noteservice.New(noteRepository)
	noteHandler := notedelivery.New(noteService)

	return newRouter(db, cfg, noteHandler)
}

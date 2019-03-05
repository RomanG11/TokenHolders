package main

import (
	"TokenHolders/cmd/initializer"
	"TokenHolders/internal/pkg/repo/models"
	"github.com/rs/zerolog/log"
)

func main() {
	app := initializer.InitApplication()

	err := app.Repo.DB.AutoMigrate(&models.Holder{}).Error
	if err != nil {
		log.Panic().Err(err).Msg("Migration error")
	}
}

package main

import (
	"TokenHolders/cmd/initializer"
	"TokenHolders/internal/distribution"
	"github.com/rs/zerolog/log"
)

func main() {
	app := initializer.InitApplication()

	err := distribution.Airdrop(app)
	if err != nil {
		log.Panic().Err(err).Msg("something went wrong")
	}
}

package main

import (
	"TokenHolders/cmd/initializer"
	"TokenHolders/internal/app"
	"github.com/rs/zerolog/log"
)

func main() {
	appl := initializer.InitApplication()

	app.RunDCListener(appl)

	//err := distribution.Airdrop(appl)
	//if err != nil {
	//	log.Panic().Err(err).Msg("something went wrong")
	//}

	log.Info().Msg("Airdrop completed")
}

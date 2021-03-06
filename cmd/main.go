package main

import (
	"TokenHolders/cmd/initializer"
	"TokenHolders/internal/app"
	"github.com/rs/zerolog/log"
)

func main() {

	appl := initializer.InitApplication()

	app.RunListener(appl)
	log.Info().Msg("Listener action completed")

	//go app.FinalCheck(appl)
	//
	//err := app.FinalCheckBack(appl)
	//if err != nil {
	//	log.Error().Err(err).Msg("FinalCheck error")
	//}
	//
	//log.Info().Msg("FinalCheck action completed")

	app.CheckBalances(appl)
}

package app

import (
	"TokenHolders/internal/pkg/application"
	"TokenHolders/internal/pkg/repo/models"
	"github.com/rs/zerolog/log"
)

func FinalCheck(app *application.Application) error {

	for {
		holders, err := app.Repo.Holder.FindAll()
		if err != nil {
			return err
		}

		for i, holder := range holders {
			b, err := app.Client.CheckFinalBalance(holder.EthAddress)
			if err != nil {
				log.Error().Err(err).Msgf("%v", b)
			}

			check(app, holder)

			if i%100 == 0 {
				log.Debug().Msgf("checked %d holders", i)
			}
		}

		log.Info().Msg("check completed")
		return nil
	}
}

func FinalCheckBack(app *application.Application) error {

	for {
		holders, err := app.Repo.Holder.FindAll()
		if err != nil {
			return err
		}

		for i := len(holders) - 1; i > 0; i-- {
			b, err := app.Client.CheckFinalBalance(holders[i].EthAddress)
			if err != nil {
				log.Error().Err(err).Msgf("%v", b)
			}

			check(app, holders[i])

			if i%100 == 0 {
				log.Debug().Msgf("checked %d holders", i)
			}
		}

		log.Info().Msg("check completed")
		return nil
	}
}

func check(app *application.Application, holder models.Holder) error {
	b, err := app.Client.CheckFinalBalance(holder.EthAddress)
	if err != nil {
		log.Error().Err(err).Msgf("%v", b)
		return err
	}

	if !holder.Balance.Equal(b) {
		log.Warn().Msgf("balances is not equals for account: %s. Previous: %s; current: %s. Changing balance",
			holder.EthAddress, holder.Balance.String(), b.String())

		holder.Balance = b

		err = app.Repo.Holder.UpdateHolder(&holder)
		if err != nil {
			return err
		}
	}
	return nil
}

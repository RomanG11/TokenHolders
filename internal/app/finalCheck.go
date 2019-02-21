package app

import (
	"TokenHolders/internal/pkg/application"
	"github.com/rs/zerolog/log"
)

func FinalCheck(app *application.Application) error {
	//var s int64 = 1
	//var f int64 = 101

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

			if !holder.Balance.Equal(b) {
				log.Warn().Msgf("balances is not equals for account: %s. Previous: %s; current: %s. Changing balance",
					holder.EthAddress, holder.Balance.String(), b.String())

				holders[i].Balance = b

				err = app.Repo.Holder.UpdateHolder(&holders[i])
				if err != nil {
					return err
				}
			}

			if i%100 == 0 {
				log.Debug().Msgf("checked %d holders", i)
			}
		}

		log.Info().Msg("check completed")
		return nil
	}
}

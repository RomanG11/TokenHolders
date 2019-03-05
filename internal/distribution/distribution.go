package distribution

import (
	"TokenHolders/internal/pkg/application"
	"TokenHolders/internal/pkg/etherPkg"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

//0x82b962d1d9b334296ea7af3e6a6018ffa1b36ad9
func Airdrop(app *application.Application) error {

	holders, err := app.Repo.Holder.FindAllWithPositiveBalance()
	if err != nil {
		return err
	}

	var addresses []common.Address
	var values []*big.Int

	key, err := crypto.HexToECDSA(app.Client.PrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("invalid ethereum private key")
		return err
	}

	auth := bind.NewKeyedTransactor(key)
	auth.GasPrice = big.NewInt(10000000000) //10 GWEI

	j := 0

	for i, holder := range holders {

		bal, err := app.Client.CheckFinalBalance(holder.EthAddress)
		if err != nil {
			log.Error().Err(err).Msgf("%v", bal)
			return err
		}

		if !bal.Equal(decimal.New(0, 0)) {

			if bal.GreaterThan(holder.Balance) {
				log.Warn().Msgf("holder %s have balance %s, should have: %s", holder.EthAddress, bal.String(), holder.Balance.String())
			} else if bal.Equal(holder.Balance) {
				log.Warn().Msgf("holder %s already got his tokens correctly", holder.EthAddress)
			} else {
				log.Warn().Msgf("holder %s already got his tokens, but have another amount", holder.EthAddress)
			}

			continue
		}

		if j == 0 {
			j = i
		}

		addresses = append(addresses, common.HexToAddress(holder.EthAddress))

		a := big.Int{}
		val, ok := a.SetString(holder.Balance.String(), 10)
		if !ok {
			return fmt.Errorf("invalid balance, %s", holder.Balance.String())
		}
		values = append(values, val)

		if len(addresses) == 50 || i == len(holders)-1 {

			tx := sendTx(app.Client, auth, addresses, values)
			log.Info().Msgf("transfer from %d to %d completed. TX hash: %s", j, i, tx.Hash().String())

			awaitTx(app, addresses[0].String(), addresses[25].String(), addresses[40].String())

			addresses = []common.Address{}
			values = []*big.Int{}

			j = 0
		}
	}

	return nil
}

func sendTx(cl *etherPkg.Client, auth *bind.TransactOpts, addresses []common.Address, values []*big.Int) *types.Transaction {

	tx, err := cl.Token.Airdrop(auth, addresses, values)
	if err != nil {
		log.Info().Err(err).Msg("Cannot send transaction, retrying in 10 seconds")
		time.Sleep(10 * time.Second)
		sendTx(cl, auth, addresses, values)
	}

	return tx
}

func awaitTx(app *application.Application, a, b, c string) {

	bal, err := app.Client.CheckFinalBalance(a)
	if err != nil {
		log.Error().Err(err).Msgf("%v", bal)
	}

	if bal.Equal(decimal.New(0, 0)) {
		time.Sleep(5 * time.Second)
		log.Info().Msg("Not completed yet")

		awaitTx(app, a, b, c)
		return
	}

	bal, err = app.Client.CheckFinalBalance(b)
	if err != nil {
		log.Error().Err(err).Msgf("%v", bal)
	}

	if bal.Equal(decimal.New(0, 0)) {
		time.Sleep(5 * time.Second)
		log.Info().Msg("Not completed yet")

		awaitTx(app, a, b, c)
		return
	}

	bal, err = app.Client.CheckFinalBalance(c)
	if err != nil {
		log.Error().Err(err).Msgf("%v", bal)
	}

	if bal.Equal(decimal.New(0, 0)) {
		time.Sleep(5 * time.Second)
		log.Info().Msg("Not completed yet")

		awaitTx(app, a, b, c)
		return
	}

	log.Info().Msg("check completed")
}

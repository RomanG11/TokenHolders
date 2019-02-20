package distribution

import (
	"TokenHolders/internal/pkg/application"
	"TokenHolders/internal/pkg/etherPkg"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"math/big"
)

func Airdrop(app *application.Application) error {
	var s int64 = 1
	var f int64 = 51

	for {
		holders, err := app.Repo.Holder.FindGroup(s, f)
		if err != nil {
			return err
		}

		var addresses []common.Address
		var values []*big.Int

		for _, holder := range holders {

			addresses = append(addresses, common.HexToAddress(holder.EthAddress))

			a := big.Int{}
			val, ok := a.SetString(holder.Balance.String(), 10)
			if !ok {
				return fmt.Errorf("invalid balance, %s", holder.Balance.String())
			}
			values = append(values, val)
		}

		key, err := crypto.HexToECDSA(app.Client.PrivateKey)
		if err != nil {
			log.Error().Err(err).Msg("invalid ethereum private key")
			return err
		}

		auth := bind.NewKeyedTransactor(key)

		sendTx(app.Client, auth, addresses, values)

		addresses = []common.Address{}
		values = []*big.Int{}

		log.Info().Msgf("transfer from %d to %d completed", s, f)

		s += 50
		f += 50

		if len(holders) != 50 {
			break
		}
	}

	return nil
}

func sendTx(cl *etherPkg.Client, auth *bind.TransactOpts, addresses []common.Address, values []*big.Int) {
	_, err := cl.Token.Airdrop(auth, addresses, values)
	if err != nil {
		log.Info().Err(err).Msg("Cannot send transaction, retrying")
		sendTx(cl, auth, addresses, values)
	}
}

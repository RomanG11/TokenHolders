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

	j := 0

	for i, holder := range holders {

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

		if len(addresses) == 50 {
			sendTx(app.Client, auth, addresses, values)

			addresses = []common.Address{}
			values = []*big.Int{}

			log.Info().Msgf("transfer from %d to %d completed", j, i)

			j = 0
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

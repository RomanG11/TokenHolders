package app

import (
	"context"
	"encoding/hex"
	"errors"
	"TokenHolders/internal/pkg/application"
	"TokenHolders/internal/pkg/repo"
	"TokenHolders/internal/pkg/repo/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"math/big"
)

const eventTransfer  = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
const iter = int64(1000)


func RunListener(app *application.Application) {

	currentBlock := app.Client.FromBlock
	lastBlock := app.Client.LastBlock
	for {
		fb := currentBlock

		currentBlock += iter
		if currentBlock > lastBlock {
			currentBlock = lastBlock
		}

		filter := ethereum.FilterQuery{
			Addresses: []common.Address{app.Client.TokenAddress},
			Topics: [][]common.Hash{{common.HexToHash(eventTransfer)}},
			FromBlock: big.NewInt(fb),
			ToBlock: big.NewInt(currentBlock),
		}

		logs, err:= app.Client.EthClient.FilterLogs(context.Background(), filter)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get filter logs for contract")
			return
		}

		for _, res := range logs {
			err := checkTransferLog(app.Repo, res)
			if err != nil {
				log.Error().Err(err).Msgf("%v", res)
			}
		}

		log.Debug().Msgf("listened to %d block", currentBlock)

		if currentBlock == lastBlock {
			break
		}
	}

	log.Debug().Msg("listenEvents successfully finished")
}

func checkTransferLog(repo *repo.Repo, ethLog types.Log) error {
	var fromStr = "0x" + ethLog.Topics[1].String()[26:]
	var toStr = "0x" + ethLog.Topics[2].String()[26:]

	var value = new(big.Int)
	value, ok := value.SetString(hex.EncodeToString(ethLog.Data), 16)
	if !ok {
		log.Debug().Msg("cannot parse uint from value")
		return errors.New("cannot parse uint from value")
	}

	log.Debug().Msgf("New Transfer event detected. From: %s, To: %s, value: %v", fromStr, toStr, value)


	zero := decimal.New(0, 0)
	var from models.Holder
	from, err := repo.Holder.GetHolderByAddress(fromStr)
	if err != nil {
		from, err = repo.Holder.NewHolder(fromStr, zero)
		if err != nil {
			return err
		}
	}

	var to models.Holder
	to, err = repo.Holder.GetHolderByAddress(toStr)
	if err != nil {
		to, err = repo.Holder.NewHolder(toStr,zero)
		if err != nil {
			return err
		}
	}

	from.Balance = from.Balance.Sub(decimal.NewFromBigInt(value, 0))
	to.Balance = to.Balance.Add(decimal.NewFromBigInt(value, 0))

	err = repo.Holder.UpdateHolder(&from)
	if err != nil {
		return err
	}

	err = repo.Holder.UpdateHolder(&to)
	if err != nil {
		return err
	}

	return nil
}

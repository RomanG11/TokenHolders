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

func RunListener(app *application.Application) {

	var filter ethereum.FilterQuery

	filter.Addresses = []common.Address{app.Client.TokenAddress}
	filter.Topics = [][]common.Hash{{common.HexToHash(eventTransfer)}}
	filter.ToBlock = big.NewInt(app.Client.LastBlock)
	logs, filterErr := app.Client.EthClient.FilterLogs(context.Background(), filter)

	if filterErr != nil {
		log.Error().Err(filterErr).Msgf("failed to get filter logs for contract ")
		return
	}

	for _, res := range logs {
		err := checkTransferLog(app.Repo, res)
		if err != nil {
			log.Error().Err(err).Msgf("%v", res)
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

	var from models.Holder
	from, err := repo.Holder.GetHolderByAddress(fromStr)
	if err != nil {
		from, err = repo.Holder.NewHolder(fromStr, decimal.New(0, 0))
		if err != nil {
			return err
		}
	}

	var to models.Holder
	to, err = repo.Holder.GetHolderByAddress(toStr)
	if err != nil {
		to, err = repo.Holder.NewHolder(toStr, decimal.New(0, 0))
	}

	from.Balance = from.Balance.Sub(decimal.NewFromBigInt(value, 0))
	to.Balance = from.Balance.Add(decimal.NewFromBigInt(value, 0))

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

package app

import (
	"TokenHolders/cmd/initializer"
	"TokenHolders/internal/pkg/application"
	"TokenHolders/internal/pkg/repo"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"math/big"
)

const iterDC = int64(100)

func RunDCListener(app *application.Application) {

	currentBlock := int64(7283241)
	lastBlock := int64(7283941)

	for {
		fb := currentBlock + 1
		currentBlock += iterDC

		if currentBlock > lastBlock {
			currentBlock = lastBlock
		}

		filter := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress("0xc690F7C7FcfFA6a82b79faB7508c466FEfdfc8c5")},
			FromBlock: big.NewInt(fb),
			ToBlock:   big.NewInt(currentBlock),
		}

	loop:
		logs, err := app.Client.EthClient.FilterLogs(context.Background(), filter)
		if err != nil {
			a := initializer.InitApplication()
			app.Client = a.Client
			log.Error().Err(err).Msgf("failed to get filter logs for contract")
			goto loop
			return
		}

		for _, res := range logs {
			if res.Topics[0] == common.HexToHash(eventTransfer) {
				err := checkDCTransferLog(app.Repo, res)
				if err != nil {
					log.Error().Err(err).Msgf("%v", res)
				}
			}
		}

		log.Debug().Msgf("listened to %d block", currentBlock)

		if currentBlock == lastBlock {
			break
		}
	}

	log.Debug().Msg("listenEvents successfully finished")
}

func checkDCTransferLog(repo *repo.Repo, ethLog types.Log) error {
	var fromStr = "0x" + ethLog.Topics[1].String()[26:]
	var toStr = "0x" + ethLog.Topics[2].String()[26:]

	if fromStr == "0xadeac69acaa3fc48416e9350436a3dd72ab0d30d" {
		h, err := repo.Holder.FindHolder(toStr)
		if err != nil {
			return err
		}

		h.Ok = true

		err = repo.Holder.UpdateHolder(&h)
		if err != nil {
			return err
		}
	}

	return nil
}

package app

import (
	"TokenHolders/internal/pkg/application"
	"TokenHolders/internal/pkg/etherPkg/contracts/tokenContract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

func finalCheck(app *application.Application, address string) (decimal.Decimal, error) {
	var d decimal.Decimal

	token, err := tokenContract.NewToken(app.Client.TokenAddress, app.Client.EthClient)
	if err != nil {
		return d, err
	}

	co := bind.CallOpts{BlockNumber: big.NewInt(app.Client.LastBlock)}

	b, err := token.BalanceOf(&co, common.HexToAddress(address))
	if err != nil {
		return d, err
	}

	return decimal.NewFromBigInt(b, 0), nil
}
package etherPkg

import (
	"TokenHolders/internal/pkg/etherPkg/contracts/tokenContract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"time"
)

type Client struct {
	EthClient *ethclient.Client
	TokenAddress common.Address
	FromBlock int64
	LastBlock int64
	Token *tokenContract.Token

	PrivateKey string
}

func InitClient(rpcPort, tokenAddress, fromBlock, lastBlock, privateKey string) *Client {
	var cl Client
	cl.EthClient = getEthClient(rpcPort)
	cl.TokenAddress = common.HexToAddress(tokenAddress)

	fb, err := strconv.ParseInt(fromBlock, 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid lastBlock")
	}

	cl.FromBlock = fb

	lb, err := strconv.ParseInt(lastBlock, 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid lastBlock")
	}

	cl.LastBlock = lb

	token, err := tokenContract.NewToken(cl.TokenAddress, cl.EthClient)
	if err != nil {
		log.Fatal().Err(err).Msg("create token instance error")
	}

	cl.Token = token
	cl.PrivateKey = privateKey

	return &cl
}

//GetEthClient initializing and return ethereum rpc client
func getEthClient(rpcPort string) *ethclient.Client {
	log.Info().Msg("Setting WS connection")
	ethClient, err := ethclient.Dial(rpcPort)
	if err != nil {
		log.Error().Err(err).Msg("Node connection broken. Reconnecting in 5 second")
		time.Sleep(5000 * time.Millisecond)
		getEthClient(rpcPort)
	}

	return ethClient
}

func (c *Client) CheckFinalBalance(address string) (decimal.Decimal, error) {
	var d decimal.Decimal

	co := bind.CallOpts{BlockNumber: big.NewInt(c.LastBlock)}

	b, err := c.Token.BalanceOf(&co, common.HexToAddress(address))
	if err != nil {
		return d, err
	}

	return decimal.NewFromBigInt(b, 0), nil
}
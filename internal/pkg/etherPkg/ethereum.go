package etherPkg

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type Client struct {
	EthClient *ethclient.Client
	TokenAddress common.Address
	LastBlock int64
}

func InitClient(rpcPort, tokenAddress, lastBlock string) *Client {
	var cl Client
	cl.EthClient = getEthClient(rpcPort)
	cl.TokenAddress = common.HexToAddress(tokenAddress)
	fd, err := strconv.ParseInt(lastBlock, 10, 64)
	if err != nil {
		log.Info().Err(err).Msg("LastBlock setting to latest")
	} else {
		cl.LastBlock = fd
	}

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

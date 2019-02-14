package initializer

import (
	"TokenHolders/internal/pkg/application"
	"TokenHolders/internal/pkg/etherPkg"
	"TokenHolders/internal/pkg/repo"
	"TokenHolders/pkg/enviroment"
	"github.com/rs/zerolog/log"
)

func InitApplication() *application.Application {
	var app application.Application
	env := enviroment.GetEnv()

	app.Repo = initDB(env)
	app.Client = initClient(env)

	return &app
}

func initDB(env map[string]string) *repo.Repo {
	name := env["DATABASE_NAME"]
	user := env["DATABASE_USER"]
	pass := env["DATABASE_PASSWORD"]
	host := env["DATABASE_HOST"]
	port := env["DATABASE_PORT"]

	r, err := repo.GetDbClient(name, user, pass, host, port)
	if err != nil {
		log.Panic().Err(err).Msg("cannot connect to database")
	}

	return r
}

func initClient(env map[string]string) *etherPkg.Client {
	rpcPort := env["RPC_PORT"]
	tokenAddress := env["TOKEN_ADDRESS"]
	lastBlock := env["LAST_BLOCK"]

	return etherPkg.InitClient(rpcPort, tokenAddress, lastBlock)
}
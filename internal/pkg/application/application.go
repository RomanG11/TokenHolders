package application

import (
	"TokenHolders/internal/pkg/etherPkg"
	"TokenHolders/internal/pkg/repo"
)

type Application struct {
	Repo                 *repo.Repo
	Client				 *etherPkg.Client
}

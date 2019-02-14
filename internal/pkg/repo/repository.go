package repo

import (
	"github.com/jinzhu/gorm"
	//postgres database driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Repo contains all table repositories
type Repo struct {
	DB                      *gorm.DB
	Holder                  HolderRepository
}

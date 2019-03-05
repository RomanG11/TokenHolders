package repo

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

//GetDbClient initializing new database client
func GetDbClient(name, user, pass, host, port string) (*Repo, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, pass, name, host, port))

	db.LogMode(false)
	if err != nil {
		return nil, err
	}

	return &Repo{
		db,
		&HolderRepo{db},
		&HolderNewRepo{db},
	}, nil
}

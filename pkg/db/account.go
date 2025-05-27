package db

import (
	"context"
	"log"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/models"
	"github.com/Bank/pkg/utils"
)

// Remove DB everywhere
type DB struct {
	db string
}

type Account interface {
	InsertAccountInDB(models.Account) error
	DeleteAcountInDB(string) error
	SearchAccountInDB(string) (models.Account, error)
	UpdateAccount(string, models.Account) (models.Account, error)
}

func NewAccount() Account {
	return &DB{}
}

func (d *DB) InsertAccountInDB(input models.Account) error {
	log.Println("Inside Account Function ")
	_, err := utils.Connection.NewInsert().Model(&input).Exec(context.Background())

	if err != nil {
		log.Println("Error in inserting in db", err)
		return err
	}
	log.Println("Successfully Inserted into DB")
	return nil
}

func (d *DB) DeleteAcountInDB(input string) error {
	var acnt models.Account
	rowsAffected, err := utils.Connection.NewDelete().Model(&acnt).Where("customer_id=?", input).Exec(context.Background())

	if err != nil {
		log.Println("Error in deleting field", err)
		return err
	}

	result, err := rowsAffected.RowsAffected()

	if result == 0 || err != nil {
		log.Println("No such record found", err)
		return customerrors.ErrNotFound
	}

	return nil

}

func (d *DB) SearchAccountInDB(input string) (models.Account, error) {
	var acnt models.Account
	err := utils.Connection.NewSelect().Model(&acnt).Where("customer_id=?", input).Scan(context.Background())

	if err != nil {
		log.Println("Error in searching field", err)
		return acnt, err
	}

	return acnt, nil
}

func (d *DB) UpdateAccount(input string, account models.Account) (models.Account, error) {
	responseDb, err := utils.Connection.NewUpdate().Model(&account).Where("customer_id=?", input).Exec(context.Background())
	if err != nil {
		log.Println("Error in updating record.", err)
		return models.Account{}, err
	}
	log.Println("Verifying rows affected")
	rowsAffected, err := responseDb.RowsAffected()
	if rowsAffected == 0 || err != nil {
		log.Println("No such record found", err)
		return models.Account{}, customerrors.ErrNotFound
	}
	log.Println("Succesfully updated record.")
	return account, nil
}

package moneyhandler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/db"
	"github.com/Bank/pkg/models"
	"github.com/gorilla/mux"
)

func AddDeposit(w http.ResponseWriter, r *http.Request) {
	// Create a map to store incoming request
	var m map[string]float32

	// Decode the incoming request
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.Println("Error from reading request.")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// Validate Input
	if m["Balance"] <= 0 {
		log.Println("Invalid Balance Deposit Request")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// Extract Customer ID from path parameters
	customerId := mux.Vars(r)

	// Call the DB function
	dbFn := db.NewTransaction()
	updatedBalance, err := dbFn.AddMoney(customerId["uid"], m["Balance"])

	// Check for errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Send appropriate response
	fResponse := fmt.Sprintf("Money Deposited Succesfully. Your updated Balance is %.2f", updatedBalance)
	w.Write([]byte(fResponse))
}

func WithdrawMoney(w http.ResponseWriter, r *http.Request) {

	// Create a map to store incoming request
	var m map[string]float32

	// Decode the incoming request
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.Println("Error in decoding request.")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// Validate Input
	if m["Balance"] <= 0 {
		log.Println("Invalid Balance Deposit Request")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// Extract Customer ID from path parameters
	customerId := mux.Vars(r)

	// Call the DB function
	dbFn := db.NewTransaction()
	remainingBalance, err := dbFn.WithdrawMoney(customerId["uid"], m["Balance"])

	// Check for errors
	if err != nil {
		if errors.Is(err, customerrors.ErrInvalidInput) {
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		} else if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, customerrors.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	// Send appropriate response
	response := fmt.Sprintf("Your remaining Balance is : %.2f", remainingBalance)

	w.Write([]byte(response))
}

func TransferMoney(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside Transfer Money")

	var trsfr models.TransferMoney
	err := json.NewDecoder(r.Body).Decode(&trsfr)
	if err != nil {
		log.Println("Error from reading request.")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Input Rec:", trsfr)

	if trsfr.Amount <= 0 {
		log.Println("Invalid Amount")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Input verified:")

	senderId := mux.Vars(r)
	// Create a Transaction
	log.Println("Creating db interface")

	t := db.NewTransaction()
	log.Println("Calling db fn")
	remianingBalance, err := t.TransferMoney(trsfr, senderId["uid"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusNotFound)
			return
		} else if errors.Is(err, customerrors.ErrInvalidInput) {
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		} else if errors.Is(err, customerrors.ErrNotFound) {
			http.Error(w, customerrors.ErrNotFound.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}
	}
	response := fmt.Sprintf("Money transferred successfully. Your remaining Balance is:%.2f", remianingBalance)
	w.Write([]byte(response))
}

package accountv1handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/db"
	"github.com/Bank/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Accounts interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	DeleteAccount(w http.ResponseWriter, r *http.Request)
	GetAccount(w http.ResponseWriter, r *http.Request)
}

type accounts struct {
	db db.Account
}

func NewAccount(ac db.Account) Accounts {
	return &accounts{}
}

func (a *accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc models.Account

	// Read Incoming Request , check for error (if any)
	err := json.NewDecoder(r.Body).Decode(&acc)

	if err != nil {
		log.Println("Error in reading request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Set defaults
	acc.CustomerID = uuid.New()
	acc.Status = "Active"

	// Check if any essential field is empty

	if acc.Balance < 0 {
		log.Println("Invalid Balance")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	if acc.Type == "" {
		log.Println("Account type not defined.")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// INSERT INTO DB & CHECK FOR ANY ERROR

	err = a.db.InsertAccountInDB(acc)
	// a.db.
	if err != nil {
		log.Println("Data Insertion Failed", err)
		http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Succesfully created new Account."))
}

func (a *accounts) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Store the deleting criteria received through path parameter
	mapOfPathPara := mux.Vars(r)

	log.Println("Value of path parameter:", mapOfPathPara)
	_, err := uuid.Parse(mapOfPathPara["uid"])
	if err != nil {
		log.Println("Invalid UID received.")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// Create dbInterface variable of db interface so that we can access underlying methods
	dbInterface := db.NewAccount()

	// Call the Delete Account method and pass the parameter & check errors,if any

	err = dbInterface.DeleteAcountInDB(mapOfPathPara["uid"])

	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			log.Println("Error in deleting record", err)
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("Error in deleting record", err)
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Deleted Account Successfully"))
}

func (a *accounts) GetAccount(w http.ResponseWriter, r *http.Request) {
	// Store the searching criteria received through path parameter
	mapOfPathPara := mux.Vars(r)

	log.Println("Value of path parameter:", mapOfPathPara)

	// Create a variable of db interface so that we can access underlying methods
	dbInterface := db.NewAccount()

	// Call the Delete Account method and pass the parameter & check errors,if any
	response, err := dbInterface.SearchAccountInDB(mapOfPathPara["uid"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Invalid searching criteria.", err)
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("Error in searching field", err)
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

	}

	finalResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error in marshalling data", err)
		http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(finalResponse)
}

func (a *accounts) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var newAcc models.Account
	err := json.NewDecoder(r.Body).Decode(&newAcc)
	if err != nil {
		log.Println("Error in reading from request.")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	searchingCriteria := mux.Vars(r)
	log.Println("Searching criterial:", searchingCriteria)

	db := db.NewAccount()
	oldAccount, err := db.SearchAccountInDB(searchingCriteria["uid"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No such record found.")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("Error in updating record")
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}
	}

	if newAcc.Type != "" {
		oldAccount.Type = newAcc.Type
	}

	if newAcc.Status != "" {
		oldAccount.Status = newAcc.Status
	}

	updatedAcc, err := db.UpdateAccount(searchingCriteria["uid"], oldAccount)
	if err != nil {
		log.Println("Failed to update the requested field:", err)
		http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(updatedAcc)
	if err != nil {
		log.Println("Failed to marshal data:", err)
		http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

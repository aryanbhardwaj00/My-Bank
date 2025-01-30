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

func CreateAccount(w http.ResponseWriter, r *http.Request) {
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

	if acc.Type == "" {
		log.Println("Account type not defined.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Account type not defined."))
		return
	}

	// INSERT INTO DB & CHECK FOR ANY ERROR

	dbInterface := db.NewAccount()
	err = dbInterface.InsertAccountInDB(acc)
	if err != nil {
		log.Println("Data Insertion Failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Data Insertion Failed."))
		return
	}

	w.Write([]byte("Succesfully created new Account."))
}

func DeleteAcount(w http.ResponseWriter, r *http.Request) {
	// Store the deleting criteria received through path parameter
	mapOfPathPara := mux.Vars(r)

	log.Println("Value of path parameter:", mapOfPathPara)

	// Create dbInterface variable of db interface so that we can access underlying methods
	dbInterface := db.NewAccount()

	// Call the Delete Account method and pass the parameter & check errors,if any

	err := dbInterface.DeleteAcountInDB(mapOfPathPara["uid"])

	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			log.Println("Error in deleting record", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No such record found."))
			return
		} else {
			log.Println("Error in deleting record", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to delete record."))
			return
		}
	}
	w.Write([]byte("Deleted Account Successfully"))
}

func SearchAccount(w http.ResponseWriter, r *http.Request) {
	// Store the searching criteria received through path parameter
	mapOfPathPara := mux.Vars(r)

	log.Println("Value of path parameter:", mapOfPathPara)

	// Create a variable of db interface so that we can access underlying methods
	dbInterface := db.NewAccount()

	// Call the Delete Account method and pass the parameter & check errors,if any
	response, err := dbInterface.SearchAccountInDB(mapOfPathPara["uid"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Error in searching field", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No such record found."))
			return
		} else {
			log.Println("Error in searching field", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to search the requested record."))
			return
		}

	}

	finalResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error in marshalling data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(finalResponse)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var newAcc models.Account
	err := json.NewDecoder(r.Body).Decode(&newAcc)
	if err != nil {
		log.Println("Error in reading from request.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in reading from request."))
		return
	}

	searchingCriteria := mux.Vars(r)
	log.Println("Searching criterial:", searchingCriteria)

	db := db.NewAccount()
	oldAccount, err := db.SearchAccountInDB(searchingCriteria["uid"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No such record found.")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No record found."))
			return
		} else {
			log.Println("Error in updating record")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to update the requested field."))
			return
		}
	}

	if newAcc.Balance != 0 {
		if newAcc.Balance < 0 {
			log.Println("Invalid balance")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Balance."))
			return
		}
		oldAccount.Balance = newAcc.Balance
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Update failed."))
		return
	}

	response, err := json.Marshal(updatedAcc)
	if err != nil {
		log.Println("Failed to marshal data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshal data."))
	}

	w.Write(response)
}

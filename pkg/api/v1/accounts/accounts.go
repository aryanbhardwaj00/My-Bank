package accountv1handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Bank/pkg/db"
	"github.com/Bank/pkg/models"
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

	// INSERT INTO DB & CHECK FOR ANY ERROR

	dbInterface := db.NewAccount()
	err = dbInterface.InsertAccountInDB(acc)
	if err != nil {
		log.Println("Data Insertion Failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		responseMsg := "Data Insertion Failed"
		w.Write([]byte(responseMsg))
		return
	}

	resonseMsg := "Succesfully created new Account"
	w.Write([]byte(resonseMsg))
}

func DeleteAcount(w http.ResponseWriter, r *http.Request) {
	// Store the deleting criteria received through path parameter
	mapOfPathPara := mux.Vars(r)

	log.Println("Value of path parameter:", mapOfPathPara)

	// Create dbInterface variable of db interface so that we can access underlying methods
	dbInterface := db.NewAccount()

	// Call the Delete Account method and pass the parameter & check errors,if any

	err := dbInterface.DeleteAcountInDB(mapOfPathPara["acc_id"])

	if err != nil {
		log.Println("Error in deleting record", err)
		w.WriteHeader(http.StatusBadRequest)
		Response := "Failed to delete data , bad request"
		w.Write([]byte(Response))
		return
	}

	Response := "Deleted Account Successfully"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Response))
}

func SearchAccount(w http.ResponseWriter, r *http.Request) {
	// Store the searching criteria received through path parameter
	mapOfPathPara := mux.Vars(r)

	log.Println("Value of path parameter:", mapOfPathPara)

	// Create a variable of db interface so that we can access underlying methods
	dbInterface := db.NewAccount()

	// Call the Delete Account method and pass the parameter & check errors,if any
	response, err := dbInterface.SearchAccountInDB(mapOfPathPara["acc_id"])
	if err != nil {
		log.Println("Error in searching field", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	finalResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error in marshalling data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(finalResponse)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request)  {
	
}
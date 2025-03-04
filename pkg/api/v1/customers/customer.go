package customerv1handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/db"
	"github.com/Bank/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var cust models.Customer
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cust)
	// json.Decoder , decodes/ reads from request body(r.body) and
	// Decode(&cust) will store whatever was read from request body and store it in "cust" variable

	if err != nil {
		log.Println("Error in reading from request", err)
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// Set the UID and status
	cust.UID = uuid.New()
	cust.Status = "Active"

	// Check if the essential fields are empty or not
	if cust.Name == "" {
		log.Println("Empty Name")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	if cust.Age <= 0 {
		log.Println("Invalid Age")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	if cust.Contact == 0 || len(strconv.Itoa(cust.Contact)) != 10 {
		log.Println("Invalid contact number")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	isValidEmail := func(email string) bool {
		_, err := mail.ParseAddress(email)
		return err == nil
	}

	if isValidEmail(cust.PrimaryEmail) == false {
		log.Println("Invalid Primary Email")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	if cust.SecondaryEmail != "" {
		if isValidEmail(cust.SecondaryEmail) == false {
			log.Println("Invalid Secondary Email")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
	}

	if cust.Address.City == "" || cust.Address.State == "" {
		log.Println("Empty Address")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	// Creating an Instance of Customer Interface, so that we can use its underlying methods
	c := db.NewCustomer()

	// Insert into Database or return error(if any)
	err = c.InsertIntoDB(cust)
	if err != nil {
		log.Println("Database insert failed", err)
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Created new field."))
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside Delete function in handler")
	// Extract the path parameter[UID]
	mapOfPathParameters := mux.Vars(r)
	log.Println("Map of path parameter:", mapOfPathParameters)

	c := db.NewCustomer()

	err := c.DeleteCustomerInDB(mapOfPathParameters["uid"])

	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			log.Println("No such record found.")
			http.Error(w, customerrors.ErrNotFound.Error(), http.StatusNotFound)
			return
		} else {
			log.Println("Error in deleting record", err)
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Deleted the requested field."))
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	// Extract value from path parameters
	// Call the SearchCustomer function
	// Return if the requested field not found
	// If found , send it as response
	log.Println("Inside Search handler")
	m := mux.Vars(r)
	log.Println("Path parameters:", m)
	result, err := db.NewCustomer().GetCustomerInDB(m["uid"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, customerrors.ErrNotFound.Error(), http.StatusNotFound)
			return
		} else {
			log.Println("Cannot find the requested field", err)
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}
	}
	log.Println("marshalling data")
	response, err := json.Marshal(result)
	if err != nil {
		log.Println("Unable to send", err)
		http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {

	var updtCust models.Customer

	// Decode and store the incoming data and check for errors

	err := json.NewDecoder(r.Body).Decode(&updtCust)
	if err != nil {
		log.Println("error in reading request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Store the searching criteria(received through path parameter) in a variable
	searchingCriteria := mux.Vars(r)

	// First fetch the record which needs to be updated

	dbInterface := db.NewCustomer()
	oldCust, err := dbInterface.GetCustomerInDB(searchingCriteria["uid"])
	log.Println("Fetched data from DB:", oldCust)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Error in searching for the requested field:", err)
			http.Error(w, customerrors.ErrNotFound.Error(), http.StatusNotFound)
			return
		} else {
			log.Println("Error in searching: ", err)
			http.Error(w, customerrors.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

	}

	// Validate each field and make updates in fetched record

	if updtCust.Name != oldCust.Name {
		if updtCust.Name == "" {
			log.Println("Empty name field.")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
		oldCust.Name = updtCust.Name
	}

	if updtCust.Age != oldCust.Age {
		if updtCust.Age <= 0 {
			log.Println("Invalid value received for age")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
		oldCust.Age = updtCust.Age
		log.Println("updated age", oldCust)
	}

	if updtCust.Contact != oldCust.Contact {
		if len(strconv.Itoa(updtCust.Contact)) == 10 && updtCust.Contact > 0 {
			oldCust.Contact = updtCust.Contact
			log.Println("updated contact", oldCust)
		} else {
			log.Println("Invalid contact")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
	}

	// Validate Email
	isValidEmail := func(email string) bool {
		_, err := mail.ParseAddress(email)
		return err == nil
	}

	if updtCust.PrimaryEmail != oldCust.PrimaryEmail {
		if isValidEmail(updtCust.PrimaryEmail) == true {
			oldCust.PrimaryEmail = updtCust.PrimaryEmail
			log.Println("updated primary email", oldCust)
		} else {
			log.Println("Invalid Primary Email")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
	}

	if updtCust.SecondaryEmail != "" {
		if updtCust.SecondaryEmail != oldCust.SecondaryEmail {
			if isValidEmail(updtCust.SecondaryEmail) == true {
				oldCust.SecondaryEmail = updtCust.SecondaryEmail
				log.Println("updated secondary email", oldCust)
			} else {
				log.Println("Invalid Email")
				http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
				return
			}
		}

	}

	if updtCust.Status != oldCust.Status {
		if updtCust.Status == "" {
			log.Println("Invalid status")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
		oldCust.Status = updtCust.Status
		log.Println("updated status", oldCust)
	}

	if updtCust.Address.City != oldCust.Address.City {
		if updtCust.Address.City == "" {
			log.Println("Invalid city")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
		oldCust.Address.City = updtCust.Address.City
	}

	if updtCust.Address.State != oldCust.Address.State {
		if updtCust.Address.State == "" {
			log.Println("Invalid State")
			http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}
		oldCust.Address.State = updtCust.Address.State
	}

	// After making all the necessary changes call the UpdateCustomer function

	updatedRecord, err := dbInterface.UpdateCustomerInDB(searchingCriteria["uid"], oldCust)

	if err != nil {
		log.Println("Error in updating requested field", err)
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusInternalServerError)
		return
	}
	
	// Marshal and send back response
	finalResponse, err := json.Marshal(updatedRecord)
	if err != nil {
		log.Println("Error in marshalling", err)
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(finalResponse))
}

func ListCustomer(w http.ResponseWriter, r *http.Request) {
	var listPara models.Listparameters
	// Unmarshal incoming data
	err := json.NewDecoder(r.Body).Decode(&listPara)

	// Convert necessary fields to lowercase for proper validation
	order := strings.ToLower(listPara.OrderBy)
	sort := strings.ToLower(listPara.SortIn)
	listPara.SortIn = sort
	listPara.OrderBy = order

	if listPara.PageSize == 0 {
		listPara.PageSize = 3
	}

	// Validate the inputs received
	validation := models.ValidateListParam(listPara)
	if validation == false {
		log.Println("Validation Failed")
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	arr := models.GetListColumns()
	input := fmt.Sprintf("order by Case when %v='%v' then 0 else 1 End , name %v offset %v  limit %v ", listPara.OrderBy, listPara.Input, listPara.SortIn, listPara.PageNumber, listPara.PageSize)
	// Call the List Customer funciton
	result, err := db.NewCustomer().ListCustomers(arr, input)
	if err != nil {
		log.Println("Cannot find the requested field", err)
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusInternalServerError)
		return
	}
	// Send back response
	response, err := json.Marshal(result)
	if err != nil {
		log.Println("Unable to send", err)
		http.Error(w, customerrors.ErrInvalidInput.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

package moneyhandler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type AddTestCase struct {
	name               string
	input              uuid.UUID
	inputBody          map[string]float32
	expectedError      string
	expectedStatusCode int
}

// Create Test Cases
var AddTest = []AddTestCase{
	{
		name:               "Valid Inputs",
		input:              uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody:          map[string]float32{"Balance": 1000},
		expectedError:      "",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "Negative Balance",
		input:              uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody:          map[string]float32{"Balance": -1000},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:               "Non Existing Account",
		input:              uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4887b"),
		inputBody:          map[string]float32{"Balance": 1000},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusNotFound,
	},
}

func TestDepositMoney(t *testing.T) {
	// Run a loop over each test case
	for _, tt := range AddTest {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the data
			reqBody, err := json.Marshal(tt.inputBody)
			if err != nil {
				t.Fatalf("Error in marshalling :%v", err)
			}

			// Create a  request
			req, err := http.NewRequest("POST", "http://localhost:8080/services/transaction1/"+tt.input.String(), bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Error in generating request %v", err)
			}

			// Create a client to make the request
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error in HTTP request %v", err)
			}

			// Read the response to check for errors message sent by our API
			errorMsg, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error during HTTP response reading: %v ", err)
			}
			defer res.Body.Close()

			// Asser  Status Codes and errors
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}

type WithdrawTestCase struct {
	name               string
	input              uuid.UUID
	inputBody          map[string]float32
	expectedError      string
	expectedStatusCode int
}

var WithdrawTest = []WithdrawTestCase{
	{
		name:               "Valid Inputs",
		input:              uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody:          map[string]float32{"Balance": 10},
		expectedError:      "",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "Negative Balance",
		input:              uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody:          map[string]float32{"Balance": -1000},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:               "Insufficient Balance",
		input:              uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody:          map[string]float32{"Balance": 100000},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:               "Non-Existing Account",
		input:              uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886c"),
		inputBody:          map[string]float32{"Balance": 10},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusNotFound,
	},
}

func TestWithdrawMoney(t *testing.T) {
	// Run a loop over each test case
	for _, tt := range WithdrawTest {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the data
			reqBody, err := json.Marshal(tt.inputBody)
			if err != nil {
				t.Fatalf("Error in marshalling :%v", err)
			}

			// Create a  request
			req, err := http.NewRequest("POST", "http://localhost:8080/services/transaction2/"+tt.input.String(), bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Error in generating request %v", err)
			}

			// Create a client to make the request
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error in HTTP request %v", err)
			}

			// Read the response to check for errors message sent by our API
			errorMsg, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error during HTTP response reading: %v ", err)
			}
			defer res.Body.Close()

			// Asser  Status Codes and errors
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}

type TransferTest struct {
	name               string
	input              uuid.UUID
	inputBody          models.TransferMoney
	expectedError      string
	expectedStatusCode int
}

var TransferCase = []TransferTest{
	{
		name:  "Valid Inputs",
		input: uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody: models.TransferMoney{
			Amount:     500,
			ReceiverID: "9d102b87-e3eb-4b04-8547-8fa70d681479",
		},
		expectedError:      "",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:  "Insuffiecient Balance",
		input: uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody: models.TransferMoney{
			Amount:     500000,
			ReceiverID: "9d102b87-e3eb-4b04-8547-8fa70d681479",
		},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:  "Invalid Balance",
		input: uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886b"),
		inputBody: models.TransferMoney{
			Amount:     -500,
			ReceiverID: "9d102b87-e3eb-4b04-8547-8fa70d681479",
		},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:  "Non Existing Sender ",
		input: uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886c"),
		inputBody: models.TransferMoney{
			Amount:     500,
			ReceiverID: "9d102b87-e3eb-4b04-8547-8fa70d681479",
		},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusNotFound,
	},
	{
		name:  "Non Existing Receiver ",
		input: uuid.MustParse("34ede6bf-0ddc-4d1a-aa48-b6f46da4886c"),
		inputBody: models.TransferMoney{
			Amount:     500,
			ReceiverID: "9d102b87-e3eb-4b04-8547-8fa70d681478",
		},
		expectedError:      customerrors.ErrInvalidInput.Error(),
		expectedStatusCode: http.StatusNotFound,
	},
}

func TestTransferMoney(t *testing.T) {
	// Run a loop over each test case
	for _, tt := range TransferCase {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the data
			reqBody, err := json.Marshal(tt.inputBody)
			if err != nil {
				t.Fatalf("Error in marshalling :%v", err)
			}

			// Create a  request
			req, err := http.NewRequest("POST", "http://localhost:8080/services/transaction3/"+tt.input.String(), bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Error in generating request %v", err)
			}

			// Create a client to make the request
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error in HTTP request %v", err)
			}

			// Read the response to check for errors message sent by our API
			errorMsg, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error during HTTP response reading: %v ", err)
			}
			defer res.Body.Close()

			// Asser  Status Codes and errors
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}

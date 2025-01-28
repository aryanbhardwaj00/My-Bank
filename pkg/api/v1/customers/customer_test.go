package customerv1handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Bank/pkg/customerrors"
	apperrors "github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type createTestcase struct {
	name               string
	inp                models.Customer
	expectedStatusCode int
	expectedError      error
}

var createTest = []createTestcase{
	{
		name: "Valid input",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      nil,
	},
	{
		name: "Invalid Name",
		inp: models.Customer{
			Name:           "",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Empty Age field",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            0,
			Contact:        9845647219,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Invalid Age",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            -18,
			Contact:        9845647219,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Empty Contact",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        0,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Invalid Contact-Less than ten digits",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        123456789,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Negative value for Contact",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        -1234567891,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Invalid Primary email",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        1234567891,
			PrimaryEmail:   "nishant123",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Invalid Secondary email",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        1234567891,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Empty AccountID",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        1234567891,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      0,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Negative AccountID",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        1234567891,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      -123,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "Empty Address field",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        1234567891,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      456,
			Address:        &models.Address{},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "No City in Address",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        1234567891,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      123,
			Address:        &models.Address{State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
	{
		name: "No State in Address",
		inp: models.Customer{
			Name:           "Nishant",
			Age:            23,
			Contact:        1234567891,
			PrimaryEmail:   "nishant@gmail.com",
			SecondaryEmail: "nishant1@gmail.com",
			AccountID:      100,
			Address:        &models.Address{City: "Jaipur"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput,
	},
}

func TestCreateCustomer(t *testing.T) {
	for _, tt := range createTest {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the data which is to be sent in request Body
			reqBody, err := json.Marshal(tt.inp)
			if err != nil {
				t.Fatalf("Error during marshalling:%v", err)
			}
			// Create a request
			req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/customers", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			// Create a client and make a request
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error during HTTP request: %v ", err)
			}
			// Ensure the response body is closed after reading
			defer res.Body.Close()

			// Assert Status Codes
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}
}

type srchTestcase struct {
	name               string
	inp                uuid.UUID
	expectedStatusCode int
	expectedError      error
}

var searchTest = []srchTestcase{
	{
		name:               "With Valid Inputs",
		inp:                uuid.MustParse("33343231-3837-6566-2d64-6234642d3436"),
		expectedStatusCode: http.StatusOK,
		expectedError:      nil,
	},
	{
		name:               "With non existing element",
		inp:                uuid.MustParse("33343231-3837-6566-2d64-6234642d3437"),
		expectedStatusCode: http.StatusNotFound,
		expectedError:      apperrors.ErrNotFound,
	},
}

func TestSearchCustomer(t *testing.T) {
	for _, tt := range searchTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/customers/"+tt.inp.String(), nil)
			assert.NoError(t, err)

			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error during HTTP request: %v", err)
			}
			defer res.Body.Close()

			// Assert status code
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)

		})
	}
}

type dltTestcase struct {
	name               string
	inp                uuid.UUID
	expectedStatusCode int
	expectedError      error
}

var deleteTest = []dltTestcase{
	{
		name:               "With Valid Inputs",
		inp:                uuid.MustParse("33343231-3837-6566-2d64-6234642d3436"),
		expectedStatusCode: http.StatusOK,
		expectedError:      nil,
	},
	{
		name:               "With non existing element",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3436"),
		expectedStatusCode: http.StatusNotFound,
		expectedError:      apperrors.ErrNotFound,
	},
}

func TestDel(t *testing.T) {
	for _, tt := range deleteTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/customers/"+tt.inp.String(), nil)
			assert.NoError(t, err)

			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error during HTTP request: %v", err)
			}
			// Ensure the response body is closed after reading
			defer res.Body.Close()

			// Assert status code
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)

		})
	}
}

type updtTestCase struct {
	name               string
	inp                uuid.UUID
	inpBody            models.Customer
	expectedStatusCode int
	expectedError      error
}

var updtTest = []updtTestCase{
	{
		name: "Update with valid UID",
		inp:  uuid.MustParse("30396165-6263-3965-2d34-3338322d3434"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      nil,
	},
	{
		name: "Update with Non Existing UID",
		inp:  uuid.MustParse("30396165-6263-3965-2d34-3338322d3425"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      130,
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusNotFound,
		expectedError:      customerrors.ErrNotFound,
	},

	{
		name:               "Update Age with Negative input",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3434"),
		inpBody:            models.Customer{Age: -25},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput,
	},

	{
		name:               "Update Contact with Negative input",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3434"),
		inpBody:            models.Customer{Contact: -9825471296},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput,
	},

	{
		name:               "Update Contact with Short input",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3434"),
		inpBody:            models.Customer{Contact: 982547},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput,
	},

	{
		name:               "Update Primary Email with Invalid format",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3434"),
		inpBody:            models.Customer{PrimaryEmail: "nishant123@"},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput,
	},

	{
		name:               "Update Secondary Email with Invalid format",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3434"),
		inpBody:            models.Customer{SecondaryEmail: "nishant123@"},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput,
	},

	{
		name:               "Update AccountID with negative input",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3434"),
		inpBody:            models.Customer{AccountID: -123},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput,
	},
}

func TestUpdate(t *testing.T) {
	for _, tt := range updtTest {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.inpBody)
			if err != nil {
				t.Fatalf("Error in marshalling :%v", err)
			}
			req, err := http.NewRequest("PATCH", "http://localhost:8080/api/v1/customers/"+tt.inp.String(), bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("Error in creating request: %v", err)
			}
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error during HTTP request: %v", err)
			}
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}
}

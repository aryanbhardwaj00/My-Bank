package customerv1handler_test

import (
	"bytes"
	"encoding/json"
	"io"
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
	expectedError      string
}

var createTest = []createTestcase{
	{
		name: "Valid input",
		inp: models.Customer{
			Name:           "Ritesh",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "ritesh@gmail.com",
			SecondaryEmail: "rishi1@gmail.com",
			AccountID:      129,
			Address:        &models.Address{City: "Jaipur", State: "Rajasthan"},
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "Created new field.",
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
		expectedError:      customerrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
	},
	{
		name: "Invalid Contact->Less than ten digits",
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
			req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/customer", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			// Create a client and make a request
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error during HTTP request: %v ", err)
			}
			errorMsg, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error during HTTP response reading: %v ", err)
			}
			// Ensure the response body is closed after reading
			defer res.Body.Close()

			// Assert Status Codes
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}

type dltTestcase struct {
	name               string
	inp                uuid.UUID
	expectedStatusCode int
	expectedError      string
}

var deleteTest = []dltTestcase{
	{
		name:               "With Valid Inputs",
		inp:                uuid.MustParse("33343231-3837-6566-2d64-6234642d3436"), // Already Deleted
		expectedStatusCode: http.StatusOK,
		expectedError:      "Deleted the requested field.",
	},
	{
		name:               "With non existing element",
		inp:                uuid.MustParse("30396165-6263-3965-2d34-3338322d3436"),
		expectedStatusCode: http.StatusNotFound,
		expectedError:      apperrors.ErrNotFound.Error(),
	},
	{
		name:               "Invalid UID",
		inp:                uuid.MustParse("30396165-6263-3965-2d34"),
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      apperrors.ErrInvalidInput.Error(),
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
			errorMsg, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error during HTTP response reading: %v ", err)
			}
			// Ensure the response body is closed after reading
			defer res.Body.Close()

			// Assert status code
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}

type updtTestCase struct {
	name               string
	inp                uuid.UUID
	inpBody            models.Customer
	expectedStatusCode int
	expectedError      string
}

var updtTest = []updtTestCase{
	{
		name: "Update with valid UID",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
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
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusNotFound,
		expectedError:      customerrors.ErrNotFound.Error(),
	},

	{
		name: "Update Age with Negative input",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            -22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},
	{
		name: "Update Age with No input",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            0,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "Update Contact with Negative input",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        -9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "Update Contact with Short input",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        98456,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "Update Primary Email with Invalid format",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "Update Secondary Email with Invalid format",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234",
			AccountID:      100,
			Status:         "Active",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "Update Status with Invalid format",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Status:         "",
			Address:        &models.Address{City: "Hisar", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "Update City with Invalid format",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Address:        &models.Address{City: "", State: "Haryana"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "Update State with Invalid format",
		inp:  uuid.MustParse("f8b63a48-67e2-4635-aaf1-13438dd2b7db"),
		inpBody: models.Customer{
			Name:           "Nishant",
			Age:            22,
			Contact:        9845647219,
			PrimaryEmail:   "godara@gmail.com",
			SecondaryEmail: "nishant1234@gmail.com",
			AccountID:      100,
			Address:        &models.Address{City: "Hisar", State: ""},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
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
			errorMsg, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error during HTTP response reading: %v ", err)
			}
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}

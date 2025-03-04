package customerv1handler

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

type getTestcase struct {
	name               string
	inp                uuid.UUID
	expectedStatusCode int
	expectedError      string
}

var getTest = []getTestcase{
	{
		name:               "With Valid Inputs",
		inp:                uuid.MustParse("078c4418-ebd9-44b0-bf9b-52d2a0f0ccb0"),
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{
		name:               "With non existing element",
		inp:                uuid.MustParse("33343231-3837-6566-2d64-6234642d3437"),
		expectedStatusCode: http.StatusNotFound,
		expectedError:      apperrors.ErrNotFound.Error(),
	},
}

func TestGetCustomer(t *testing.T) {
	for _, tt := range getTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/customers/"+tt.inp.String(), nil)
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
			defer res.Body.Close()

			// Assert status code
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}

type listTestCase struct {
	name               string
	input              models.Listparameters
	expectedStatusCode int
	expectedError      string
}

var listCases = []listTestCase{
	{
		name: "With Valid Inputs",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "name",
			SortIn:     "asc",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{
		name: "With Empty Input",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "",
			OrderBy:    "Name",
			SortIn:     "Asc",
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},
	{
		name: "With Empty Order By",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "",
			SortIn:     "Asc",
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},
	{
		name: "With Invalid Order By",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "Sakshi",
			SortIn:     "Asc",
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},

	{
		name: "With Valid Order By[UpperCase]",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "NAME",
			SortIn:     "Asc",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{
		name: "With Valid Order By[LowerCase]",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "name",
			SortIn:     "Asc",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{

		name: "With Valid Order By[Mixed]",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "NamE",
			SortIn:     "Asc",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{
		name: "With Empty Sort In",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "Name",
			SortIn:     "",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{
		name: "With Valid Sort In[UpperCase]",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "NAME",
			SortIn:     "ASC",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{
		name: "With Valid Sort In[LowerCase]",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "name",
			SortIn:     "asc",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{

		name: "With Valid Sort In[Mixed]",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "name",
			SortIn:     "aSc",
		},
		expectedStatusCode: http.StatusOK,
		expectedError:      "",
	},
	{
		name: "With Invalid Sort By",
		input: models.Listparameters{
			PageNumber: 0,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "Name",
			SortIn:     "acs",
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},
	{
		name: "With Invalid Page Number",
		input: models.Listparameters{
			PageNumber: -1,
			PageSize:   3,
			Input:      "Sakshi",
			OrderBy:    "Name",
			SortIn:     "asc",
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},
	{
		name: "With Invalid Page Size",
		input: models.Listparameters{
			PageNumber: 1,
			PageSize:   -3,
			Input:      "Sakshi",
			OrderBy:    "Name",
			SortIn:     "asc",
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedError:      customerrors.ErrInvalidInput.Error(),
	},
}

func TestListCustomer(t *testing.T) {
	for _, tt := range listCases {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.input)

			if err != nil {
				t.Fatalf("Error in marshalling :%v", err)
			}

			req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/customers", bytes.NewBuffer(jsonData))
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
			defer res.Body.Close()

			// Assert status code
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Contains(t, string(errorMsg), tt.expectedError)
		})
	}
}
package accountv1handler

import (
	"testing"

	"github.com/Bank/pkg/utils"
	"github.com/google/uuid"
)

type testscase struct {
	uid                uuid.UUID
	Type               string
	Balance            int
	Status             string
	CustomerID         int
	expectedStatusCode int
	expectedError      string
}

var tests = []testscase{
	{
		uid:                uuid.New(),
		Type:               "Savings",
		Balance:            0,
		Status:             "Active",
		CustomerID:         0,
		expectedStatusCode: 200,
		expectedError:      "nil",
	},
	{
		uid:        uuid.New(),
		Type:       "Current",
		Balance:    0.0,
		Status:     "Closed",
		CustomerID: 1,
	},
}

func TestCreateAccount(t *testing.T) {
	// Connect to DB
	utils.ConnectToDB()

	// Run a loop through each test case

	// for _, tt:=range tests{
	//  t.Run(tt.)
	// }
	// Check for results
}

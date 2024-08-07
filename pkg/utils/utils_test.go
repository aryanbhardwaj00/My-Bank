package utils

import (
	"log"
	"testing"
)

func TestConnectToDb(t *testing.T) {
	err := ConnectToDB()
	if err != nil {
		log.Println("Test failed")
	}
}

package main

import (
	"log"

	"github.com/KidMuon/unbearable_traffic/data_import"
)

func main() {
	err := data_import.ImportOverpassData_Standard()
	if err != nil {
		log.Fatal(err)
	}
}

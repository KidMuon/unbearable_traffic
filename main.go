package main

import (
	"github.com/KidMuon/unbearable_traffic/data_import"
)

func main() {
	err := data_import.TestPostOverpassAPI()
	if err != nil {
		panic(err)
	}
}

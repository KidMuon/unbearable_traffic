package data_import

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

func TestPostOverpassAPI() error {
	body := `
		way ["highway"~"."]["maxspeed"~"."]
		(33.435, -112.08, 33.445, -112.07);
		(._;>;);
		out body;
	`

	resp, err := http.Post(overpassAPIURL, "text/plain;charset=UTF-8", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	extractedWays := OverpassStreetResponse{}
	derr := xml.NewDecoder(resp.Body).Decode(&extractedWays)
	if derr != nil {
		return derr
	}

	fmt.Println(len(extractedWays.Streets))

	return nil
}

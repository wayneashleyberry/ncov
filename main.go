package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/leekchan/accounting"
	"github.com/wayneashleyberry/truecolor/pkg/color"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	c := &http.Client{
		Timeout: time.Second * 10,
	}

	const baseURL = "https://exchange.vcoud.com/coronavirus/latest"

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return err
	}

	// req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code")
	}

	type items []struct {
		ID             string    `json:"_id"`
		Name           string    `json:"name"`
		NameEs         string    `json:"nameEs"`
		TotalCases     int       `json:"totalCases"`
		TotalDeaths    int       `json:"totalDeaths"`
		SeriousCases   int       `json:"seriousCases"`
		TotalRecovered int       `json:"totalRecovered"`
		TotalCases24H  int       `json:"totalCases24h"`
		TotalDeaths24H int       `json:"totalDeaths24h"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
		Slug           string    `json:"slug,omitempty"`
		Symbol         string    `json:"symbol,omitempty"`
	}

	var r items

	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}

	for _, item := range r {
		if item.Name == "Total" {
			color.White().Underline().Print("Total")
			fmt.Print("  ")
			color.Color(0, 176, 101).Printf("Confirmed Cases: %s", accounting.FormatNumber(item.TotalCases, 0, ",", "."))
			fmt.Print("  ")
			color.Color(214, 49, 68).Printf("Deceased: %s", accounting.FormatNumber(item.TotalDeaths, 0, ",", "."))
			fmt.Print("  ")
			color.Color(68, 155, 226).Printf("Recovered: %s", accounting.FormatNumber(item.TotalRecovered, 0, ",", "."))
			fmt.Print("  ")
			color.Color(175, 176, 62).Printf("Serious: %s", accounting.FormatNumber(item.SeriousCases, 0, ",", "."))
			fmt.Println("")
		}

		if item.Name == "United Kingdom" {
			color.White().Underline().Print("UK")
			fmt.Print("  ")
			color.Color(0, 176, 101).Printf("Confirmed Cases: %s", accounting.FormatNumber(item.TotalCases, 0, ",", "."))
			fmt.Print("  ")
			color.Color(214, 49, 68).Printf("Deceased: %s", accounting.FormatNumber(item.TotalDeaths, 0, ",", "."))
			fmt.Print("  ")
			color.Color(68, 155, 226).Printf("Recovered: %s", accounting.FormatNumber(item.TotalRecovered, 0, ",", "."))
			fmt.Print("  ")
			color.Color(175, 176, 62).Printf("Serious: %s", accounting.FormatNumber(item.SeriousCases, 0, ",", "."))
			fmt.Println("")
		}
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
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

	names := []string{
		"Total",
		"United Kingdom",
		"South Africa",
		// "China",
		// "Australia",
		"USA",
	}

	var t time.Time

	for _, item := range r {
		for _, name := range names {
			if item.Name == name {
				color.White().Underline().Print(item.Name)
				fmt.Print("  ")
				color.Color(4, 173, 151).Printf("Confirmed: %s", accounting.FormatNumber(item.TotalCases, 0, ",", "."))
				fmt.Print("  ")
				color.Color(236, 57, 44).Printf("Deceased: %s", accounting.FormatNumber(item.TotalDeaths, 0, ",", "."))
				fmt.Print("  ")
				color.Color(52, 152, 219).Printf("Recovered: %s", accounting.FormatNumber(item.TotalRecovered, 0, ",", "."))
				fmt.Print("  ")
				color.Color(243, 156, 17).Printf("Serious: %s", accounting.FormatNumber(item.SeriousCases, 0, ",", "."))
				fmt.Print("\n")
			}
		}

		t = item.UpdatedAt
	}

	fmt.Printf("Updated %s\n", humanize.Time(t))

	return nil
}

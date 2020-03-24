package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/leekchan/accounting"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:  "ncov",
		Long: "SARS-CoV-2 / COVID-19 statistics from https://covid19stats.live",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return printStatistics(args[0])
		},
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type item struct {
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

func getItems() ([]item, error) {
	c := &http.Client{
		Timeout: time.Second * 10,
	}

	const baseURL = "https://exchange.vcoud.com/coronavirus/latest"

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return []item{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return []item{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []item{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []item{}, fmt.Errorf("bad status code")
	}

	var r []item

	err = json.Unmarshal(body, &r)
	if err != nil {
		return []item{}, err
	}

	return r, nil
}

func printStatistics(name string) error {
	items, err := getItems()
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.Name == name {
			fmt.Println(accounting.FormatNumber(item.TotalCases, 0, ",", "."))
		}
	}

	return nil
}

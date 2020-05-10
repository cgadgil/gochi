package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FieldToQueryMap() (queryMap map[string]string) {
	queryMap = make(map[string]string)
	queryMap["Cash-To-Debt"] = "#financial-strength > div:nth-child(1) > table:nth-child(3) > tbody:nth-child(2) > tr:nth-child(1) > td:nth-child(2)"
	queryMap["Altman Z-Score"] = "#financial-strength > div > table.stock-indicator-table > tbody > tr:nth-child(7) > td:nth-child(2)"
	queryMap["Operating Margin %"] = "#profitability > div > table.stock-indicator-table > tbody > tr:nth-child(1) > td:nth-child(2)"
	queryMap["Net Margin %"] = "#profitability > div > table.stock-indicator-table > tbody > tr:nth-child(2) > td:nth-child(2)"
	queryMap["PE Ratio"] = "#ratios > div > table.stock-indicator-table > tbody > tr:nth-child(1) > td:nth-child(2)"
	queryMap["Forward PE Ratio"] = "#ratios > div > table.stock-indicator-table > tbody > tr:nth-child(2) > td:nth-child(2)"
	queryMap["PEG Ratio"] = "#ratios > div > table.stock-indicator-table > tbody > tr:nth-child(12) > td:nth-child(2)"
	queryMap["Company Name"] = ".fs-x-large"
	queryMap["Current Stock Price"] = ".fs-x-large"
	queryMap["Industry"] = "#business-description > div > div:nth-child(2) > a:nth-child(4)"
	return
}

func GetValues(doc *goquery.Document) (values map[string]string) {
	values = make(map[string]string)
	for key, value := range FieldToQueryMap() {
		s := doc.Find(value).First()
		fieldValue := strings.TrimSpace(s.Text())
		if key == "Company Name" {
			fieldValue = strings.TrimSpace(strings.Split(fieldValue, "$")[0])
			//fmt.Println("Splitting...", fieldValue)
		} else if key == "Current Stock Price" {
			fieldValue = strings.TrimSpace(strings.Split(fieldValue, "$")[1])
		}
		fmt.Printf("Name [%s] Value [%s]\n", key, fieldValue)
		values[key] = fieldValue
	}
	return
}

func ScrapeFinancialData(symbol string) {
	// Request the HTML page.
	res, err := http.Get("https://www.gurufocus.com/stock/INTC/summary")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	GetValues(doc)
}

func main() {
	ScrapeFinancialData("INTC")
}

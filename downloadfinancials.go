package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	queryMap["Industry"] = "#business-description > div > div:nth-child(2) > a:nth-child(3)"
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
			fieldValue = strings.TrimSpace(strings.Split(fieldValue, " ")[0])
		}
		fmt.Printf("Name [%s] Value [%s]\n", key, fieldValue)
		values[key] = fieldValue
	}
	return
}

type FinancialSummary struct {
	createdTime     *time.Time `gorm:"primary_key"`
	symbol          string     `gorm:"primary_key"`
	companyName     string
	cashToDebt      float64
	altmanZScore    float64
	operatingMargin float64
	netMargin       float64
	peRatio         float64
	forwardPeRatio  float64
	pegRatio        float64
	stockPrice      float64
	industry        string
}

// func GetValuesFinancialSummary(doc *goquery.Document) *FinancialSummary {
// 	v := GetValues(doc)
// 	//fs := FinancialSummary { symbol: }
// }

func ScrapeFinancialData(createdTime *time.Time, symbol string) *FinancialSummary {
	// Request the HTML page.
	res, err := http.Get("https://www.gurufocus.com/stock/" + symbol + "/summary")
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

	financialData := GetValues(doc)

	myParseFloat := func(s string) float64 {
		v, _ := strconv.ParseFloat(s, 32)
		return v
	}
	return &FinancialSummary{
		createdTime:     createdTime,
		symbol:          symbol,
		companyName:     financialData["Company Name"],
		cashToDebt:      myParseFloat(financialData["Cash-To-Debt"]),
		altmanZScore:    myParseFloat(financialData["Altman Z-Score"]),
		operatingMargin: myParseFloat(financialData["Altman Z-Score"]),
		netMargin:       myParseFloat(financialData["Altman Z-Score"]),
		peRatio:         myParseFloat(financialData["Altman Z-Score"]),
		forwardPeRatio:  myParseFloat(financialData["Altman Z-Score"]),
		pegRatio:        myParseFloat(financialData["Altman Z-Score"]),
		stockPrice:      myParseFloat(financialData["Altman Z-Score"]),
		industry:        financialData["Company Name"],
	}
}

func main() {
	//db, err := gorm.Open("sqlite3", "financial_data.db")
	now := time.Now()
	ScrapeFinancialData(&now, "INTC")
}

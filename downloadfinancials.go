package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/rs/xid"
	"flag"
	"go.uber.org/ratelimit"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	queryMap["PB Ratio"] = "#ratios > div > table.stock-indicator-table > tbody > tr:nth-child(5) > td:nth-child(2)"
	queryMap["Market Cap"] = "#components-root > div > div.responsive-section > div:nth-child(1) > div > div.summary-section.summary-section-sm > div > div.stock-summary-table.fc-regular > div:nth-child(6)"
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
		log.Printf("Name [%s] Value [%s]\n", key, fieldValue)
		values[key] = fieldValue
	}
	return
}

type FinancialSummary struct {
	gorm.Model
	UniqueId string `gorm:"primary_key"`
	Symbol          string 
	CreatedTime     *time.Time
	CompanyName     string
	CashToDebt      float64 `gorm:"not null;type:decimal(10,2)"`
	AltmanZScore    float64 `gorm:"not null;type:decimal(10,2)"`
	OperatingMargin float64 `gorm:"not null;type:decimal(10,2)"`
	NetMargin       float64 `gorm:"not null;type:decimal(10,2)"`
	PeRatio         float64 `gorm:"not null;type:decimal(10,2)"`
	ForwardPeRatio  float64 `gorm:"not null;type:decimal(10,2)"`
	PegRatio        float64 `gorm:"not null;type:decimal(10,2)"`
	PbRatio float64 `gorm:"not null;type:decimal(10,2)"`
	StockPrice      float64 `gorm:"not null;type:decimal(10,2)"`
	Industry        string
}

func ScrapeFinancialData(uniqueID string, createdTime *time.Time, symbol string) *FinancialSummary {
	// Request the HTML page.
	res, err := http.Get("https://www.gurufocus.com/stock/" + symbol + "/summary")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil;
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
		UniqueId: uniqueID,
		CreatedTime:     createdTime,
		Symbol:          symbol,
		CompanyName:     financialData["Company Name"],
		CashToDebt:      myParseFloat(financialData["Cash-To-Debt"]),
		AltmanZScore:    myParseFloat(financialData["Altman Z-Score"]),
		OperatingMargin: myParseFloat(financialData["Operating Margin %"])/100,
		NetMargin:       myParseFloat(financialData["Net Margin %"])/100,
		PeRatio:         myParseFloat(financialData["PE Ratio"]),
		ForwardPeRatio:  myParseFloat(financialData["Forward PE Ratio"]),
		PegRatio:        myParseFloat(financialData["PEG Ratio"]),
		PbRatio: myParseFloat(financialData["PB Ratio"]),
		StockPrice:      myParseFloat(financialData["Current Stock Price"]),
		Industry:        financialData["Industry"],
	}
}

func main() {
	db, err := gorm.Open("sqlite3", "financial_data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Migrate the schema
	db.AutoMigrate(&FinancialSummary{})

	now := time.Now()
	guid := xid.New()

	boolPtr := flag.Bool("save", true, "save to database")
	flag.Parse()

	rl := ratelimit.New(2) // limit to 2 reqs/sec
	for _, symbol := range flag.Args() {
		rl.Take()
		financialDataSummary := ScrapeFinancialData(guid.String(), &now, symbol)
		if *boolPtr {
			if(financialDataSummary != nil) {
				db.Create(financialDataSummary)
			}
		} else {
			fmt.Println(financialDataSummary)
		}
	}


}

package sheets

import (
	// "fmt"
	// "io"
	"encoding/json"
	"net/http"
	"log"
	"bytes"
	// "errors"
	"strconv"
	// "strings"

	"golang.org/x/oauth2"
)

func InsertRows(sheetID string, colValues []string, tabName string) {
	conf := GoogleAuth()
	client := conf.Client(oauth2.NoContext)
	
	// POST request
    url := "https://sheets.googleapis.com/v4/spreadsheets/" +sheetID+ "/values/" +tabName+ "!A1%3AA2:append?includeValuesInResponse=false&insertDataOption=INSERT_ROWS&valueInputOption=USER_ENTERED&prettyPrint=true"

	arrayDynamic := [][]string{
		colValues, // colValues = {"H", "I", "J"}
	}
	values := map[string][][]string{"values": arrayDynamic}
	jsonValue, _ := json.Marshal(values)

    // var jsonStr = []byte(`{"values":[["A","B"]]}`) // also functionnal
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
}

// BatchUpdate
func CreateNewTab(sheetID string, tabName string, rowC int, colC int) string {
	// spreadsheets BatchUpdate is different of spreadsheets.values BatchUpdate
	// https://developers.google.com/sheets/api/guides/batchupdate?hl=en
	conf := GoogleAuth()
	client := conf.Client(oauth2.NoContext)
	
    url := "https://sheets.googleapis.com/v4/spreadsheets/"+ sheetID + ":batchUpdate"
	jsonStr := `{
		"requests": [
			{
				"addSheet": {
					"properties": {
						"title": "`+tabName+`",
						"sheetType": "GRID",
						"gridProperties": {
							"rowCount": `+strconv.Itoa(rowC)+`,
							"columnCount": `+strconv.Itoa(colC)+`
						}
					}
				}
			}
		]
	}`
	jsonByte := []byte(jsonStr)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // fmt.Println("response Headers:", resp.Header)
    // body, _ := io.ReadAll(resp.Body)
    // fmt.Println("response Body:", string(body))
	// fmt.Println("response Status:", resp.Status) // 400 Bad Request
	return resp.Status
}

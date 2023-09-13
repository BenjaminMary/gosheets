package sheets

import (
	// "fmt"
	"io"
	"encoding/json"
	"net/http"
	"log"
	// "bytes"
	"errors"
	// "strconv"
	// "strings"

	"golang.org/x/oauth2"
)

type Gsheets struct {
	SpreadsheetID string `json:"spreadsheetId"`
	Properties    struct {
		Title         string `json:"title"`
		Locale        string `json:"locale"`
		AutoRecalc    string `json:"autoRecalc"`
		TimeZone      string `json:"timeZone"`
		DefaultFormat struct {
			BackgroundColor struct {
				Red   int `json:"red"`
				Green int `json:"green"`
				Blue  int `json:"blue"`
			} `json:"backgroundColor"`
			Padding struct {
				Top    int `json:"top"`
				Right  int `json:"right"`
				Bottom int `json:"bottom"`
				Left   int `json:"left"`
			} `json:"padding"`
			VerticalAlignment string `json:"verticalAlignment"`
			WrapStrategy      string `json:"wrapStrategy"`
			TextFormat        struct {
				ForegroundColor struct {
				} `json:"foregroundColor"`
				FontFamily           string `json:"fontFamily"`
				FontSize             int    `json:"fontSize"`
				Bold                 bool   `json:"bold"`
				Italic               bool   `json:"italic"`
				Strikethrough        bool   `json:"strikethrough"`
				Underline            bool   `json:"underline"`
				ForegroundColorStyle struct {
					RgbColor struct {
					} `json:"rgbColor"`
				} `json:"foregroundColorStyle"`
			} `json:"textFormat"`
			BackgroundColorStyle struct {
				RgbColor struct {
					Red   int `json:"red"`
					Green int `json:"green"`
					Blue  int `json:"blue"`
				} `json:"rgbColor"`
			} `json:"backgroundColorStyle"`
		} `json:"defaultFormat"`
		SpreadsheetTheme struct {
			PrimaryFontFamily string `json:"primaryFontFamily"`
			ThemeColors       []struct {
				ColorType string `json:"colorType"`
				Color     struct {
					RgbColor struct {
					} `json:"rgbColor"`
				} `json:"color"`
			} `json:"themeColors"`
		} `json:"spreadsheetTheme"`
	} `json:"properties"`
	Sheets []struct {
		Properties struct {
			SheetID        int    `json:"sheetId"`
			Title          string `json:"title"`
			Index          int    `json:"index"`
			SheetType      string `json:"sheetType"`
			GridProperties struct {
				RowCount    int `json:"rowCount"`
				ColumnCount int `json:"columnCount"`
			} `json:"gridProperties"`
		} `json:"properties"`
	} `json:"sheets"`
	SpreadsheetURL string `json:"spreadsheetUrl"`
}

func sheetsError(info string) error {
	return errors.New(info)
 }

func GetSheets(sheetID string) {
	conf := GoogleAuth()
	client := conf.Client(oauth2.NoContext)

	requestURL := "https://sheets.googleapis.com/v4/spreadsheets/" + sheetID

    req, err := http.NewRequest(http.MethodGet, requestURL, nil)
    if err != nil {
        log.Fatal(err)
		sheetsError("new request error")
    }
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
		sheetsError("request error, probably wrong Spreadsheet")
    }

    bytesR, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
		sheetsError("read request response error")
    }

    // fmt.Println(string(bytesR))

	var result Gsheets
	json.Unmarshal(bytesR, &result)
	// fmt.Println(PrettyPrint(result))
	// fmt.Println(PrettyPrint(result.Sheets))
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "  ")
    return string(s)
}

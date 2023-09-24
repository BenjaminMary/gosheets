package main

import (
    "net/http"
    // "fmt"
    "time"
    "strings"
    "strconv"
    "html/template"
    "os"
    // "encoding/json"

    "example.com/sheets"

    "github.com/gin-gonic/gin"
)

// index.html
func index(c *gin.Context) {
    c.HTML(http.StatusOK, "0.index.html", "")
}

// GET cookie
func getCookieSetup(c *gin.Context) {
    // try to read if a cookie exists, return "Aucun" otherwise
    cookie, err := c.Cookie("sheetsId")
    if err != nil {
        cookie = "Aucun"
    }    
    c.HTML(http.StatusOK, "1.cookieSetup.html", gin.H{
        "Cookie": cookie,
        "ClientEmail": os.Getenv("client_email"),
    })
}
// POST cookie
func postCookieSetup(c *gin.Context) {
    // name (string): The name of the cookie to be set.
    // value (string): The value of the cookie.
    // maxAge (int): The maximum age of the cookie in seconds. If set to 0, the cookie will be deleted immediately. If set to a negative value, the cookie will be a session cookie and will be deleted when the browser is closed.
    // path (string): The URL path for which the cookie is valid. Defaults to "/", meaning the cookie is valid for all URLs.
    // domain (string): The domain for which the cookie is valid. Defaults to the current domain with "".
    // secure (bool): If set to true, the cookie will only be sent over secure (HTTPS) connections.
    // httpOnly (bool): If set to true, the cookie will be inaccessible to JavaScript and can only be sent over HTTP(S) connections.
    
    // set a cookie
    sheetId := c.PostForm("sheetId")
    cookieDurationStr := c.PostForm("cookieDuration")
    var cookieDurationInt64 int64
    var cookieDurationInt int
	cookieDurationInt64, err := strconv.ParseInt(cookieDurationStr, 10, 0)
    if err != nil {
        c.String(http.StatusBadRequest, "cookieDurationStr conversion en int64 KO: %v", err)
        return
    }	
    cookieDurationInt = int(cookieDurationInt64)

    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("sheetsId", "", -1, "/", "", false, true)
    c.SetCookie("sheetsId", sheetId, cookieDurationInt, "/", "", false, true)
    c.String(200, `<p>Cookie: %s</p><p id="hx-swap-oob1" hx-swap-oob="true">Nouveau Cookie enregistré.</p>`, sheetId)
}

// getsheets.html
func getsheets(c *gin.Context) {
    cookieSheetsId, err := c.Cookie("sheetsId")
    if err != nil {
        cookieSheetsId = "Aucun"
        c.HTML(http.StatusOK, "getsheets.html", gin.H{
            "info": "Aucun Google Sheets ID trouvé, créer un cookie avec cette information en suivant le lien ci-dessous.",
        })
        return
    }
    sheets.GetSheets(cookieSheetsId)
    c.HTML(http.StatusOK, "getsheets.html", gin.H{
        "info": "Google Sheets ID trouvé: "+cookieSheetsId,
    })
}

type ParamStruct struct {
	Param string `json:"param"`
	Valeur string `json:"valeur"`
	Info string `json:"info"`
}
// type ParamStruct struct {
// 	Param struct {
// 		Valeur string `json:"valeur"`
// 		Info string `json:"info"`
// 	} `json:"param"`
// }
// batchget.html
func batchget(c *gin.Context) {
    cookieSheetsId, err := c.Cookie("sheetsId")
    if err != nil {
        cookieSheetsId = "Aucun"
        c.HTML(http.StatusOK, "getsheets.html", gin.H{
            "info": "Aucun Google Sheets ID trouvé, créer un cookie avec cette information en suivant le lien ci-dessous.",
        })
        return
    }
    values := sheets.BatchGetSheets(cookieSheetsId, "param")
    // var paramHead, categories ParamStruct
    m := make(map[string]string)
    params := make([]ParamStruct, 2)
    i := 0
    for _, param := range values {
        params[i].Param = param[0]
        params[i].Valeur = param[1]
        params[i].Info = param[2]
        m[param[0]] = param[1]
        i += 1
    }
	// json.Unmarshal([][]byte(values), &params)
    c.HTML(http.StatusOK, "batchget.html", gin.H{
        "info": params,
        "categories": strings.Split(params[1].Valeur, ","),
        "categoriesMap": strings.Split(m["Categories"], ","),
        "map": m,
    })
}

// GET create tab.html
func getcreatetab(c *gin.Context) {
    cookieSheetsId, err := c.Cookie("sheetsId")
    if err != nil {
        cookieSheetsId = "Aucun"
        c.HTML(http.StatusOK, "2.createtab.html", gin.H{
            "info": "Aucun Google Sheets ID trouvé, créer un cookie avec cette information en suivant le lien ci-dessous.",
        })
        return
    }
    tabName := "data"
    rowC := 1
    colC := 4  
    var colValues []string
    colValues = append(colValues, "Date", "Désignation", "Catégorie", "Prix")

    var infoCreationData, infoCreationParam string
    respStatus := sheets.CreateNewTab(cookieSheetsId, tabName, rowC, colC)
    if respStatus == "200 OK"{
        sheets.InsertRows(cookieSheetsId, colValues, tabName)
        infoCreationData = "Onglet 'Data' créé."
    } else { infoCreationData = "Onglet 'Data' déjà existant." }

    tabName = "param"
    rowC = 1
    colC = 3 
    respStatus = sheets.CreateNewTab(cookieSheetsId, tabName, rowC, colC)
    if respStatus == "200 OK"{
        colValues = nil
        colValues = append(colValues, "Param", "Valeur", "Info")
        sheets.InsertRows(cookieSheetsId, colValues, tabName)
        colValues = nil
        colValues = append(colValues, "Categories", "Supermarché,Restaurant,Loisir,Véhicule", "Inscrire la liste des catégories souhaitées, l'ordre est repris.\nSéparer les noms des champs par des virgules sans espaces.")
        sheets.InsertRows(cookieSheetsId, colValues, tabName)
        infoCreationParam = "Onglet 'Param' créé."
    } else { infoCreationParam = "Onglet 'Param' déjà existant." }
    c.HTML(http.StatusOK, "2.createtab.html", gin.H{
        "info": "Google Sheets ID trouvé: "+cookieSheetsId,
        "infoCreationData": infoCreationData,
        "infoCreationParam": infoCreationParam,
    })
}

// GET InsertRows.html
func getinsertrows(c *gin.Context) {
    currentTime := time.Now()
    currentDate := currentTime.Format("2006-01-02") // YYYY-MM-DD
    cookieSheetsId, err := c.Cookie("sheetsId")
    if err != nil {
        cookieSheetsId = "Aucun"
        c.HTML(http.StatusOK, "1.cookieSetup.html", gin.H{
            "Cookie": cookieSheetsId,
            "ClientEmail": os.Getenv("client_email"),
        })
        return
    }    
    c.HTML(http.StatusOK, "3.insertrows.html", gin.H{
        "currentDate": currentDate,
    })
}

type PostInsertRows struct {
    Date string         `form:"date" binding:"required"`
    Designation string  `form:"designation" binding:"required"`
    Categorie  string   `form:"categorie" binding:"required"`
    Prix string         `form:"prix" binding:"required"`
}
// POST InsertRows.html
func postinsertrows(c *gin.Context) {
    // time.Sleep(299999999 * time.Nanosecond) // to simulate 300ms of loading in the front when submiting form
    var Form PostInsertRows
    if err := c.ShouldBind(&Form); err != nil {
        c.String(http.StatusBadRequest, "bad request: %v", err)
        return
    }
    var colValues []string
    // "=CNUM("+strings.Replace(Form.Prix, ".", ",", 1)+")")
    // fmt.Printf("before colValues, form: %#s \n", Form) // form: {2023-09-13 désig Supermarche 5.03}
    colValues = append(colValues, Form.Date, Form.Designation, Form.Categorie, 
        strings.Replace(Form.Prix, ".", ",", 1))

    cookieSheetsId, err := c.Cookie("sheetsId")
    if err != nil {
        cookieSheetsId = "Aucun"
        c.HTML(http.StatusOK, "getsheets.html", gin.H{
            "info": "Aucun Google Sheets ID trouvé, créer un cookie avec cette information en suivant le lien ci-dessous.",
        })
        return
    }
    tabName := "data"
    sheets.InsertRows(cookieSheetsId, colValues, tabName)

    tmpl := template.Must(template.ParseFiles("./html/templates/3.insertrows.html"))
    tmpl.ExecuteTemplate(c.Writer, "lastInsert", Form)
}

func main() {
    router := gin.Default()

    // render HTML
    // https://gin-gonic.com/docs/examples/html-rendering/
	router.LoadHTMLGlob("html/**/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
    
    // SERVE STATICS
    router.StaticFile("/favicon.ico", "./img/favicon.ico")
    router.StaticFile("/favicon.png", "./img/favicon.png") // 32x32
    router.Static("/img", "./img")

    router.GET("/", index)
    router.GET("/getsheets", getsheets)
    router.GET("/batchget", batchget)

    router.GET("/cookie-setup", getCookieSetup)
    router.POST("/cookie-setup", postCookieSetup)

    router.GET("/insertrows", getinsertrows)
    router.POST("/insertrows", postinsertrows)

    router.GET("/create-tab", getcreatetab)
    // router.POST("/create-tab", postcreatetab)

    router.Run("0.0.0.0:8082")
}
